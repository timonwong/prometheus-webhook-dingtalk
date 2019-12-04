// Copyright © 2016 Prometheus Team
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/google/go-github/v25/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	"github.com/prometheus/promu/util/retry"
)

var (
	releasecmd     = app.Command("release", "Upload all release files to the Github release")
	timeout        = releasecmd.Flag("timeout", "Upload timeout").Duration()
	allowedRetries = releasecmd.Flag("retry", "Number of retries to perform when upload fails").
			Default("2").Int()
	releaseLocation = releasecmd.Arg("location", "Location of files to release").Default(".").Strings()
	versionRe       = regexp.MustCompile(`^\d+\.\d+\.\d+(-.+)?$`)
)

func isPrerelease(version string) (bool, error) {
	matches := versionRe.FindStringSubmatch(version)
	if matches == nil {
		return false, errors.Errorf("invalid version %s", version)
	}
	return matches[1] != "", nil
}

func runRelease(location string) {
	token := os.Getenv("GITHUB_TOKEN")
	if len(token) == 0 {
		fatal(errors.New("GITHUB_TOKEN not defined"))
	}

	ctx := context.Background()
	if *timeout != time.Duration(0) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, *timeout)
		defer cancel()
	}
	client := github.NewClient(
		oauth2.NewClient(
			ctx,
			oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: token},
			),
		),
	)

	prerelease, err := isPrerelease(projInfo.Version)
	if err != nil {
		fatal(err)
	}

	// Find the GitHub release matching with the tag. We need to list all
	// releases because it is the only way to get draft releases too.
	var (
		release *github.RepositoryRelease
		opts    = &github.ListOptions{}
		tag     = fmt.Sprintf("v%s", projInfo.Version)
	)
	for {
		releases, resp, err := client.Repositories.ListReleases(ctx, projInfo.Owner, projInfo.Name, opts)
		if err != nil {
			fatal(errors.Wrap(err, "failed to list releases"))
		}
		for _, r := range releases {
			if r.GetTagName() == tag {
				release = r
				break
			}
		}
		if release != nil || resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	if release == nil {
		// Create a draft release if none exists already.
		var (
			err        error
			draft      = true
			name, body = getChangelog(projInfo.Version, readChangelog())
		)
		if name != "" {
			release, _, err = client.Repositories.CreateRelease(
				ctx,
				projInfo.Owner,
				projInfo.Name,
				&github.RepositoryRelease{
					TagName:         &tag,
					TargetCommitish: &projInfo.Revision,
					Name:            &name,
					Body:            &body,
					Draft:           &draft,
					Prerelease:      &prerelease,
				})
			if err != nil {
				fatal(errors.Wrap(err, fmt.Sprintf("failed to create a draft release for %s", projInfo.Version)))
			}
		} else {
			fmt.Println("fail to parse CHANGELOG.md")
		}
	}

	if err := filepath.Walk(location, releaseFile(ctx, client, release)); err != nil {
		// Remove incomplete assets.
		// See https://developer.github.com/v3/repos/releases/#response-for-upstream-failure
		opts = &github.ListOptions{}
		for {
			assets, resp, err := client.Repositories.ListReleaseAssets(ctx, projInfo.Owner, projInfo.Name, release.GetID(), opts)
			if err != nil {
				break
			}
			for _, asset := range assets {
				if strings.EqualFold(asset.GetState(), "starter") {
					_, _ = client.Repositories.DeleteReleaseAsset(ctx, projInfo.Owner, projInfo.Name, asset.GetID())
				}
			}
			if resp.NextPage == 0 {
				break
			}
			opts.Page = resp.NextPage
		}
		fatal(errors.Wrap(err, "failed to upload all files"))
	}
}

func releaseFile(ctx context.Context, client *github.Client, release *github.RepositoryRelease) func(string, os.FileInfo, error) error {
	return func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}

		// Check if the asset has already been uploaded and remove it if it is a draft release.
		filename := filepath.Base(path)
		opts := &github.ListOptions{}
		for {
			assets, resp, err := client.Repositories.ListReleaseAssets(ctx, projInfo.Owner, projInfo.Name, release.GetID(), opts)
			if err != nil {
				return errors.Wrap(err, "failed to list release assets")
			}
			var stop bool
			for _, asset := range assets {
				if asset.GetName() == filename {
					var err error
					stop = true
					if release.GetDraft() {
						_, err = client.Repositories.DeleteReleaseAsset(ctx, projInfo.Owner, projInfo.Name, asset.GetID())
						if err != nil {
							err = errors.Wrapf(err, "failed to delete existing asset %q", filename)
						}
					} else {
						err = errors.Errorf("%q already exists", filename)
					}
					if err != nil {
						return err
					}
					break
				}
			}
			if stop || resp.NextPage == 0 {
				break
			}
			opts.Page = resp.NextPage
		}

		maxAttempts := *allowedRetries + 1
		err = retry.Do(func(attempt int) (bool, error) {
			again := attempt < maxAttempts

			f, err := os.Open(path)
			if err != nil {
				return again, err
			}

			_, _, err = client.Repositories.UploadReleaseAsset(
				ctx,
				projInfo.Owner, projInfo.Name, release.GetID(),
				&github.UploadOptions{Name: filename},
				f)
			if err != nil {
				time.Sleep(2 * time.Second)
			}

			return again, err
		})
		if err != nil {
			return errors.Wrapf(err, "failed to upload %q after %d attempts", filename, maxAttempts)
		}
		fmt.Println(" > uploaded", filename)

		return nil
	}
}

func readChangelog() io.ReadCloser {
	f, err := os.Open("CHANGELOG.md")
	if err != nil {
		fmt.Printf("fail to read CHANGELOG.md: %v\n", err)
		return ioutil.NopCloser(&bytes.Buffer{})
	}
	return f
}

// getChangelog returns the changelog's header and body for a given version.
func getChangelog(version string, rc io.ReadCloser) (string, string) {
	defer rc.Close()

	var (
		scanner = bufio.NewScanner(rc)
		s       []string
		header  string
		reading bool
	)
	for (len(s) == 0 || reading) && scanner.Scan() {
		text := scanner.Text()
		switch {
		case strings.HasPrefix(text, "## "+version+" "):
			reading = true
			header = strings.TrimSpace(strings.TrimPrefix(text, "##"))
		case strings.HasPrefix(text, "## "):
			reading = false
		case reading:
			if len(s) == 0 && strings.TrimSpace(text) == "" {
				continue
			}
			s = append(s, scanner.Text())
		}
	}

	return header, strings.Join(s, "\n")
}
