// Copyright © 2017 Prometheus Team
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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

var (
	// validHeaderStrings is a slice of strings that must exist in a header.
	validHeaderStrings = []string{"copyright", "generated"}

	checkcmd         = app.Command("check", "Check the resources for validity")
	checkLicensescmd = checkcmd.Command("licenses", "Inspect source files for each file in a given directory")
	sourceExtensions = checkLicensescmd.Flag("extensions", "Comma separated list of valid source code extensions (default is .go)").
				Default(".go").Strings()
	headerLength = checkLicensescmd.Flag("length", "The number of lines to read from the head of the file").
			Short('n').Default("10").Int()
	checkLicLocation = checkLicensescmd.Arg("location", "Directory path to check licenses").
				Default(".").Strings()
)

func runCheckLicenses(path string, n int, extensions []string) {
	path = fmt.Sprintf("%s%c", filepath.Clean(path), filepath.Separator)

	filesMissingHeaders, err := checkLicenses(path, n, extensions)
	if err != nil {
		fatal(errors.Wrap(err, "Failed to check files for license header"))
	}

	for _, file := range filesMissingHeaders {
		fmt.Println(file)
	}
}

func checkLicenses(path string, n int, extensions []string) ([]string, error) {
	var missingHeaders []string
	walkFunc := func(filepath string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		if strings.HasPrefix(filepath, "vendor/") {
			return nil
		}

		if !suffixInSlice(f.Name(), extensions) {
			return nil
		}

		file, err := os.Open(filepath)
		if err != nil {
			return err
		}

		defer file.Close()

		pass := false
		scanner := bufio.NewScanner(file)
		for i := 0; i < n; i++ {
			scanner.Scan()

			if err = scanner.Err(); err != nil {
				return err
			}

			if stringContainedInSlice(strings.ToLower(scanner.Text()), validHeaderStrings) {
				pass = true
			}
		}

		if !pass {
			missingHeaders = append(missingHeaders, filepath)
		}

		return nil
	}

	err := filepath.Walk(path, walkFunc)
	if err != nil {
		return nil, err
	}

	return missingHeaders, nil
}

func stringContainedInSlice(needle string, haystack []string) bool {
	exists := false
	for _, h := range haystack {
		if strings.Contains(needle, h) {
			exists = true
			break
		}
	}

	return exists
}

func suffixInSlice(needle string, haystack []string) bool {
	exists := false
	for _, h := range haystack {
		if strings.HasSuffix(needle, h) {
			exists = true
			break
		}
	}

	return exists
}
