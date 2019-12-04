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
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/prometheus/promu/util/sh"
)

var (
	dockerBuilderImageName = "quay.io/prometheus/golang-builder"

	defaultMainPlatforms = []string{
		"linux/amd64", "linux/386", "darwin/amd64", "darwin/386", "windows/amd64", "windows/386",
		"freebsd/amd64", "freebsd/386", "openbsd/amd64", "openbsd/386", "netbsd/amd64", "netbsd/386",
		"dragonfly/amd64",
	}
	defaultARMPlatforms = []string{
		"linux/armv5", "linux/armv6", "linux/armv7", "linux/arm64", "freebsd/armv6", "freebsd/armv7",
		"openbsd/armv7", "netbsd/armv6", "netbsd/armv7",
	}
	defaultPowerPCPlatforms = []string{
		"aix/ppc64", "linux/ppc64", "linux/ppc64le",
	}
	defaultMIPSPlatforms = []string{
		"linux/mips", "linux/mipsle",
		"linux/mips64", "linux/mips64le",
	}
	defaultS390Platforms = []string{
		"linux/s390x",
	}
	armPlatformsAliases = map[string][]string{
		"linux/arm":   {"linux/armv5", "linux/armv6", "linux/armv7"},
		"freebsd/arm": {"freebsd/armv6", "freebsd/armv7"},
		"openbsd/arm": {"openbsd/armv7"},
		"netbsd/arm":  {"netbsd/armv6", "netbsd/armv7"},
	}
)

var (
	crossbuildcmd        = app.Command("crossbuild", "Crossbuild a Go project using Golang builder Docker images")
	crossBuildCgoFlagSet bool
	crossBuildCgoFlag    = crossbuildcmd.Flag("cgo", "Enable CGO using several docker images with different crossbuild toolchains.").
				PreAction(func(c *kingpin.ParseContext) error {
			crossBuildCgoFlagSet = true
			return nil
		}).Default("false").Bool()
	goFlagSet bool
	goFlag    = crossbuildcmd.Flag("go", "Golang builder version to use (e.g. 1.11)").
			PreAction(func(c *kingpin.ParseContext) error {
			goFlagSet = true
			return nil
		}).String()
	platformsFlagSet bool
	platformsFlag    = crossbuildcmd.Flag("platforms", "Space separated list of platforms to build").Short('p').
				PreAction(func(c *kingpin.ParseContext) error {
			platformsFlagSet = true
			return nil
		}).Strings()
	// kingpin doesn't currently support using the crossbuild command and the
	// crossbuild tarball subcommand at the same time, so we treat the
	// tarball subcommand as an optional arg
	tarballsSubcommand = crossbuildcmd.Arg("tarballs", "Optionally pass the string \"tarballs\" from cross-built binaries").String()
)

func runCrossbuild() {
	//Check required configuration
	if len(strings.TrimSpace(config.Repository.Path)) == 0 {
		log.Fatalf("missing required '%s' configuration", "repository.path")
	}
	if *tarballsSubcommand == "tarballs" {
		runCrossbuildTarballs()
		return
	}

	if crossBuildCgoFlagSet {
		config.Go.CGo = *crossBuildCgoFlag
	}
	if goFlagSet {
		config.Go.Version = *goFlag
	}
	if platformsFlagSet {
		config.Crossbuild.Platforms = *platformsFlag
	}

	var (
		mainPlatforms    []string
		armPlatforms     []string
		powerPCPlatforms []string
		mipsPlatforms    []string
		s390xPlatforms   []string
		unknownPlatforms []string

		cgo       = config.Go.CGo
		goVersion = config.Go.Version
		repoPath  = config.Repository.Path
		platforms = config.Crossbuild.Platforms

		dockerBaseBuilderImage    = fmt.Sprintf("%s:%s-base", dockerBuilderImageName, goVersion)
		dockerMainBuilderImage    = fmt.Sprintf("%s:%s-main", dockerBuilderImageName, goVersion)
		dockerARMBuilderImage     = fmt.Sprintf("%s:%s-arm", dockerBuilderImageName, goVersion)
		dockerPowerPCBuilderImage = fmt.Sprintf("%s:%s-powerpc", dockerBuilderImageName, goVersion)
		dockerMIPSBuilderImage    = fmt.Sprintf("%s:%s-mips", dockerBuilderImageName, goVersion)
		dockerS390XBuilderImage   = fmt.Sprintf("%s:%s-s390x", dockerBuilderImageName, goVersion)
	)

	for _, platform := range platforms {
		switch {
		case stringInSlice(platform, defaultMainPlatforms):
			mainPlatforms = append(mainPlatforms, platform)
		case stringInSlice(platform, defaultARMPlatforms):
			armPlatforms = append(armPlatforms, platform)
		case stringInSlice(platform, defaultPowerPCPlatforms):
			powerPCPlatforms = append(powerPCPlatforms, platform)
		case stringInSlice(platform, defaultMIPSPlatforms):
			mipsPlatforms = append(mipsPlatforms, platform)
		case stringInSlice(platform, defaultS390Platforms):
			s390xPlatforms = append(s390xPlatforms, platform)
		case stringInMapKeys(platform, armPlatformsAliases):
			armPlatforms = append(armPlatforms, armPlatformsAliases[platform]...)
		default:
			unknownPlatforms = append(unknownPlatforms, platform)
		}
	}

	if len(unknownPlatforms) > 0 {
		warn(errors.Errorf("unknown/unhandled platforms: %s", unknownPlatforms))
	}

	if !cgo {
		// In non-CGO, use the base image without any crossbuild toolchain
		var allPlatforms []string
		allPlatforms = append(allPlatforms, mainPlatforms[:]...)
		allPlatforms = append(allPlatforms, armPlatforms[:]...)
		allPlatforms = append(allPlatforms, powerPCPlatforms[:]...)
		allPlatforms = append(allPlatforms, mipsPlatforms[:]...)
		allPlatforms = append(allPlatforms, s390xPlatforms[:]...)

		pg := &platformGroup{"base", dockerBaseBuilderImage, allPlatforms}
		if err := pg.Build(repoPath); err != nil {
			fatal(errors.Wrapf(err, "The %s builder docker image exited unexpectedly", pg.Name))
		}
	} else {
		for _, pg := range []platformGroup{
			{"main", dockerMainBuilderImage, mainPlatforms},
			{"ARM", dockerARMBuilderImage, armPlatforms},
			{"PowerPC", dockerPowerPCBuilderImage, powerPCPlatforms},
			{"MIPS", dockerMIPSBuilderImage, mipsPlatforms},
			{"s390x", dockerS390XBuilderImage, s390xPlatforms},
		} {
			if err := pg.Build(repoPath); err != nil {
				fatal(errors.Wrapf(err, "The %s builder docker image exited unexpectedly", pg.Name))
			}
		}
	}
}

type platformGroup struct {
	Name        string
	DockerImage string
	Platforms   []string
}

func (pg platformGroup) Build(repoPath string) error {
	platformsParam := strings.Join(pg.Platforms[:], " ")
	if len(platformsParam) == 0 {
		return nil
	}

	fmt.Printf("> running the %s builder docker image\n", pg.Name)

	cwd, err := os.Getwd()
	if err != nil {
		return errors.Wrapf(err, "couldn't get current working directory")
	}

	ctrName := "promu-crossbuild-" + pg.Name + strconv.FormatInt(time.Now().Unix(), 10)
	err = sh.RunCommand("docker", "create", "-t",
		"--name", ctrName,
		pg.DockerImage,
		"-i", repoPath,
		"-p", platformsParam)
	if err != nil {
		return err
	}

	err = sh.RunCommand("docker", "cp",
		cwd+"/.",
		ctrName+":/app/")
	if err != nil {
		return err
	}

	err = sh.RunCommand("docker", "start", "-a", ctrName)
	if err != nil {
		return err
	}

	err = sh.RunCommand("docker", "cp", "-a",
		ctrName+":/app/.build/.",
		cwd+"/.build")
	if err != nil {
		return err
	}
	return sh.RunCommand("docker", "rm", "-f", ctrName)
}
