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
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/pkg/errors"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/prometheus/promu/util/sh"
)

var (
	buildcmd        = app.Command("build", "Build a Go project")
	buildCgoFlagSet bool
	buildCgoFlag    = buildcmd.Flag("cgo", "Enable CGO").
			PreAction(func(c *kingpin.ParseContext) error {
			buildCgoFlagSet = true
			return nil
		}).Bool()
	prefixFlagSet bool
	prefixFlag    = buildcmd.Flag("prefix", "Specific dir to store binaries (default is .)").
			PreAction(func(c *kingpin.ParseContext) error {
			prefixFlagSet = true
			return nil
		}).String()
	binariesArg = buildcmd.Arg("binary-names", "List of binaries to build").Default("all").Strings()
)

// Check if binary names passed to build command are in the config.
// Returns an array of Binary to build, or error.
func validateBinaryNames(binaryNames []string, cfgBinaries []Binary) ([]Binary, error) {
	var binaries []Binary

OUTER:
	for _, binaryName := range binaryNames {
		for _, binary := range cfgBinaries {
			if binaryName == binary.Name {
				binaries = append(binaries, binary)
				continue OUTER
			}
		}
		return nil, fmt.Errorf("binary %s not found in config", binaryName)
	}
	return binaries, nil
}

func buildBinary(ext string, prefix string, ldflags string, binary Binary) {
	info("Building binary: " + binary.Name)
	binaryName := fmt.Sprintf("%s%s", binary.Name, ext)
	fmt.Printf(" >   %s\n", binaryName)

	repoPath := config.Repository.Path
	flags := config.Build.Flags

	params := []string{"build",
		"-o", path.Join(prefix, binaryName),
		"-ldflags", ldflags,
	}

	params = append(params, sh.SplitParameters(flags)...)
	params = append(params, path.Join(repoPath, binary.Path))
	info("Building binary: " + "go " + strings.Join(params, " "))
	if err := sh.RunCommand("go", params...); err != nil {
		fatal(errors.Wrap(err, "command failed: "+strings.Join(params, " ")))
	}
}

func buildAll(ext string, prefix string, ldflags string, binaries []Binary) {
	for _, binary := range binaries {
		buildBinary(ext, prefix, ldflags, binary)
	}
}

func runBuild(binariesString string) {
	//Check required configuration
	if len(strings.TrimSpace(config.Repository.Path)) == 0 {
		log.Fatalf("missing required '%s' configuration", "repository.path")
	}
	if buildCgoFlagSet {
		config.Go.CGo = *buildCgoFlag
	}
	if prefixFlagSet {
		config.Build.Prefix = *prefixFlag
	}

	var (
		cgo    = config.Go.CGo
		prefix = config.Build.Prefix

		ext      string
		binaries = config.Build.Binaries
		ldflags  string
	)

	if goos == "windows" {
		ext = ".exe"
	}

	ldflags = getLdflags(projInfo)

	os.Setenv("CGO_ENABLED", "0")
	if cgo {
		os.Setenv("CGO_ENABLED", "1")
	}
	defer os.Unsetenv("CGO_ENABLED")

	if binariesString == "all" {
		buildAll(ext, prefix, ldflags, binaries)
		return
	}

	binariesArray := strings.Split(binariesString, ",")
	binariesToBuild, err := validateBinaryNames(binariesArray, binaries)
	if err != nil {
		fatal(errors.Wrap(err, "validation of given binary names for build command failed"))
	}

	for _, binary := range binariesToBuild {
		buildBinary(ext, prefix, ldflags, binary)
	}
}

func getLdflags(info ProjectInfo) string {
	var ldflags []string

	if len(strings.TrimSpace(config.Build.LDFlags)) > 0 {
		var (
			tmplOutput = new(bytes.Buffer)
			fnMap      = template.FuncMap{
				"date":     time.Now().UTC().Format,
				"host":     os.Hostname,
				"repoPath": RepoPathFunc,
				"user":     UserFunc,
			}
			ldflagsTmpl = config.Build.LDFlags
		)

		tmpl, err := template.New("ldflags").Funcs(fnMap).Parse(ldflagsTmpl)
		if err != nil {
			fatal(errors.Wrap(err, "Failed to parse ldflags text/template"))
		}

		if err := tmpl.Execute(tmplOutput, info); err != nil {
			fatal(errors.Wrap(err, "Failed to execute ldflags text/template"))
		}

		ldflags = append(ldflags, strings.Split(tmplOutput.String(), "\n")...)
	} else {
		ldflags = append(ldflags, fmt.Sprintf("-X main.Version=%s", info.Version))
	}

	extLDFlags := config.Build.ExtLDFlags
	if config.Build.Static && goos != "darwin" && goos != "solaris" && !stringInSlice("-static", extLDFlags) {
		extLDFlags = append(extLDFlags, "-static")
	}

	if len(extLDFlags) > 0 {
		ldflags = append(ldflags, fmt.Sprintf("-extldflags '%s'", strings.Join(extLDFlags, " ")))
	}

	return strings.Join(ldflags[:], " ")
}

// UserFunc returns the current username.
func UserFunc() (interface{}, error) {
	// os/user.Current() doesn't always work without CGO
	return shellOutput("whoami"), nil
}

// RepoPathFunc returns the repository path.
func RepoPathFunc() interface{} {
	return config.Repository.Path
}
