// Copyright (c) 2013 The Gocov Authors.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/axw/gocov/gocov/internal/testflag"
)

// resolvePackages returns a slice of resolved package names, given a slice of
// package names that could be relative or recursive.
func resolvePackages(pkgs []string) ([]string, error) {
	var buf bytes.Buffer
	cmd := exec.Command("go", append([]string{"list", "-e"}, pkgs...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	var resolvedPkgs []string
	lines := strings.Split(buf.String(), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			resolvedPkgs = append(resolvedPkgs, line)
		}
	}
	return resolvedPkgs, nil
}

func runTests(args []string) error {
	pkgs, testFlags := testflag.Split(args)
	pkgs, err := resolvePackages(pkgs)
	if err != nil {
		return err
	}

	tmpDir, err := ioutil.TempDir("", "gocov")
	if err != nil {
		return err
	}
	defer func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			log.Printf("failed to clean up temp directory %q", tmpDir)
		}
	}()

	// Unique -coverprofile file names are used so that all the files can be
	// later merged into a single file.
	for i, pkg := range pkgs {
		coverFile := filepath.Join(tmpDir, fmt.Sprintf("test%d.cov", i))
		cmdArgs := append([]string{"test", "-coverprofile", coverFile}, testFlags...)
		cmdArgs = append(cmdArgs, pkg)
		cmd := exec.Command("go", cmdArgs...)
		cmd.Stdin = nil
		// Write all test command output to stderr so as not to interfere with
		// the JSON coverage output.
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	// Packages without tests will not produce a coverprofile; only pick up the
	// ones that were created.
	files, err := filepath.Glob(filepath.Join(tmpDir, "test*.cov"))
	if err != nil {
		return err
	}

	// Merge the profiles.
	return convertProfiles(files...)
}
