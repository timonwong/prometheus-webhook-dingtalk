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
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/prometheus/promu/util/sh"
)

var (
	tarballcmd = app.Command("tarball", "Create a tarball from the built Go project")

	tarballPrefixSet bool
	tarballPrefix    = tarballcmd.Flag("prefix", "Specific dir to store tarballs").
				PreAction(func(c *kingpin.ParseContext) error {
			tarballPrefixSet = true
			return nil
		}).
		Default(".").String()

	tarBinariesLocation = tarballcmd.Arg("location", "location of binaries to tar").Default(".").Strings()
)

func runTarball(binariesLocation string) {
	if tarballPrefixSet {
		config.Tarball.Prefix = *tarballPrefix
	}

	var (
		prefix = config.Tarball.Prefix
		tmpDir = ".release"
		goos   = envOr("GOOS", goos)
		goarch = envOr("GOARCH", goarch)
		name   = fmt.Sprintf("%s-%s.%s-%s", projInfo.Name, projInfo.Version, goos, goarch)

		binaries = config.Build.Binaries
		ext      string
	)

	if goos == "windows" {
		ext = ".exe"
	}

	dir := filepath.Join(tmpDir, name)

	if err := os.MkdirAll(dir, 0777); err != nil {
		fatal(errors.Wrap(err, "Failed to create directory"))
	}
	defer sh.RunCommand("rm", "-rf", tmpDir)

	projectFiles := config.Tarball.Files
	for _, file := range projectFiles {
		sh.RunCommand("cp", "-a", file, dir)
	}

	for _, binary := range binaries {
		binaryName := fmt.Sprintf("%s%s", binary.Name, ext)
		sh.RunCommand("cp", "-a", filepath.Join(binariesLocation, binaryName), dir)
	}

	if !fileExists(prefix) {
		os.Mkdir(prefix, 0777)
	}

	tar := fmt.Sprintf("%s.tar.gz", name)
	fmt.Println(" >  ", tar)
	sh.RunCommand("tar", "zcf", filepath.Join(prefix, tar), "-C", tmpDir, name)
}
