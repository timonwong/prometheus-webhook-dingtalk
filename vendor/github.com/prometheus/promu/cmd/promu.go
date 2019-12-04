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
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/prometheus/promu/util/sh"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	yaml "gopkg.in/yaml.v2"
)

const (
	// DefaultConfigFilename contains the default filename of the promu config file
	DefaultConfigFilename = ".promu.yml"
)

// Binary represents a built binary.
type Binary struct {
	Name string
	Path string
}

// Config contains the Promu Command Configuration
type Config struct {
	Build struct {
		Binaries   []Binary
		Flags      string
		LDFlags    string
		ExtLDFlags []string
		Prefix     string
		Static     bool
	}
	Crossbuild struct {
		Platforms []string
	}
	Repository struct {
		Path string
	}
	Go struct {
		CGo     bool
		Version string
	}
	Tarball struct {
		Files  []string
		Prefix string
	}
}

// NewConfig creates a Config initialized with default values
// some values may be overridden by CLI args
func NewConfig() *Config {
	config := &Config{}
	config.Build.Binaries = []Binary{{Name: projInfo.Name, Path: "."}}
	config.Build.Prefix = "."
	config.Build.Static = true
	platforms := defaultMainPlatforms
	platforms = append(platforms, defaultARMPlatforms...)
	platforms = append(platforms, defaultPowerPCPlatforms...)
	platforms = append(platforms, defaultMIPSPlatforms...)
	platforms = append(platforms, defaultS390Platforms...)
	config.Crossbuild.Platforms = platforms
	config.Tarball.Prefix = "."
	config.Go.Version = "1.12"
	config.Go.CGo = false
	config.Repository.Path = projInfo.Repo

	return config
}

var (
	buildContext = build.Default
	goos         = buildContext.GOOS
	goarch       = buildContext.GOARCH

	configFile = app.Flag("config", "Path to config file").Short('c').
			Default(DefaultConfigFilename).String()
	verbose  = app.Flag("verbose", "Verbose output").Short('v').Bool()
	config   *Config
	projInfo ProjectInfo

	// app represents the base command
	app = kingpin.New("promu", "promu is the utility tool for building and releasing Prometheus projects")
)

// init prepares flags
func init() {
	app.HelpFlag.Short('h')
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	var err error
	projInfo, err = NewProjectInfo()
	checkError(err, "Unable to initialize project info")

	command := kingpin.MustParse(app.Parse(os.Args[1:]))
	sh.Verbose = *verbose
	initConfig(*configFile)

	info(fmt.Sprintf("Running command: %v %v", command, os.Args[2:]))

	switch command {
	case buildcmd.FullCommand():
		runBuild(optArg(*binariesArg, 0, "all"))
	case checkLicensescmd.FullCommand():
		runCheckLicenses(optArg(*checkLicLocation, 0, "."), *headerLength, *sourceExtensions)
	case checksumcmd.FullCommand():
		runChecksum(optArg(*checksumLocation, 0, "."))
	case crossbuildcmd.FullCommand():
		runCrossbuild()
	case infocmd.FullCommand():
		runInfo()
	case releasecmd.FullCommand():
		runRelease(optArg(*releaseLocation, 0, "."))
	case tarballcmd.FullCommand():
		runTarball(optArg(*tarBinariesLocation, 0, "."))
	case versioncmd.FullCommand():
		runVersion()
	}
}

// initConfig reads the given config file into the Config object
func initConfig(filename string) {
	info(fmt.Sprintf("Using config file: %v", filename))

	configData, err := ioutil.ReadFile(filename)
	checkError(err, "Unable to read config file: "+filename)
	config = NewConfig()
	err = yaml.Unmarshal(configData, config)
	checkError(err, "Unable to parse config file: "+filename)
}

// info prints the given message only if running in verbose mode
func info(message string) {
	if *verbose {
		fmt.Println(message)
	}
}

// warn prints a non-fatal error
func warn(err error) {
	if *verbose {
		fmt.Fprintf(os.Stderr, `/!\ %+v\n`, err)
	} else {
		fmt.Fprintln(os.Stderr, `/!\`, err)
	}
}

// printErr prints a error
func printErr(err error) {
	if *verbose {
		fmt.Fprintf(os.Stderr, "!! %+v\n", err)
	} else {
		fmt.Fprintln(os.Stderr, "!!", err)
	}
}

// fatal prints a error and exit
func fatal(err error) {
	printErr(err)
	os.Exit(1)
}

// shellOutput executes a shell command and returns the trimmed output
func shellOutput(cmd string) string {
	args := strings.Fields(cmd)
	out, _ := exec.Command(args[0], args[1:]...).Output()
	return strings.Trim(string(out), " \n\r")
}

// fileExists checks if a file exists and is not a directory
func fileExists(path ...string) bool {
	finfo, err := os.Stat(filepath.Join(path...))
	if err == nil && !finfo.IsDir() {
		return true
	}
	if os.IsNotExist(err) || finfo.IsDir() {
		return false
	}
	if err != nil {
		fatal(err)
	}
	return true
}

// readFile reads a file and return the trimmed output
func readFile(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.Trim(string(data), "\n\r ")
}

func optArg(args []string, i int, def string) string {
	if i+1 > len(args) {
		return def
	}
	return args[i]
}

func envOr(name, def string) string {
	s := os.Getenv(name)
	if s == "" {
		return def
	}
	return s
}

func stringInSlice(needle string, haystack []string) bool {
	for _, hay := range haystack {
		if hay == needle {
			return true
		}
	}
	return false
}

func stringInMapKeys(needle string, haystack map[string][]string) bool {
	_, ok := haystack[needle]
	return ok
}

// checkError prints the message and exits if the error is not nil
func checkError(e error, message string) {
	if e != nil {
		fmt.Println(message)
		fatal(e)
	}
}
