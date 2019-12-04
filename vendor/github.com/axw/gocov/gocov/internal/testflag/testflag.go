// Copyright (c) 2015 The Gocov Authors.
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
//
// Parts of this taken from cmd/go/testflag.go and
// cmd/go/build.go; adapted for simplicity.
//
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testflag

import "strings"

type testFlagSpec struct {
	name   string
	isBool bool
}

var testFlagDefn = []*testFlagSpec{
	// test-specific
	{name: "i", isBool: true},
	{name: "bench"},
	{name: "benchmem", isBool: true},
	{name: "benchtime"},
	{name: "covermode"},
	{name: "cpu"},
	{name: "cpuprofile"},
	{name: "memprofile"},
	{name: "memprofilerate"},
	{name: "blockprofile"},
	{name: "blockprofilerate"},
	{name: "parallel"},
	{name: "run"},
	{name: "short", isBool: true},
	{name: "timeout"},
	{name: "trace"},
	{name: "v", isBool: true},

	// common build flags
	{name: "a", isBool: true},
	{name: "race", isBool: true},
	{name: "x", isBool: true},
	{name: "asmflags"},
	{name: "buildmode"},
	{name: "compiler"},
	{name: "gccgoflags"},
	{name: "gcflags"},
	{name: "ldflags"},
	{name: "linkshared", isBool: true},
	{name: "pkgdir"},
	{name: "tags"},
	{name: "toolexec"},
}

// Split processes the arguments , separating flags and package
// names as done by "go test".
func Split(args []string) (packageNames, passToTest []string) {
	inPkg := false
	for i := 0; i < len(args); i++ {
		if !strings.HasPrefix(args[i], "-") {
			if !inPkg && packageNames == nil {
				// First package name we've seen.
				inPkg = true
			}
			if inPkg {
				packageNames = append(packageNames, args[i])
				continue
			}
		}

		if inPkg {
			// Found an argument beginning with "-"; end of package list.
			inPkg = false
		}

		n := parseTestFlag(args, i)
		if n == 0 {
			// This is a flag we do not know; we must assume
			// that any args we see after this might be flag
			// arguments, not package names.
			inPkg = false
			if packageNames == nil {
				// make non-nil: we have seen the empty package list
				packageNames = []string{}
			}
			passToTest = append(passToTest, args[i])
			continue
		}

		passToTest = append(passToTest, args[i:i+n]...)
		i += n - 1
	}
	return packageNames, passToTest
}

// parseTestFlag sees if argument i is a known flag and returns its
// definition, value, and whether it consumed an extra word.
func parseTestFlag(args []string, i int) (n int) {
	arg := args[i]
	if strings.HasPrefix(arg, "--") { // reduce two minuses to one
		arg = arg[1:]
	}
	switch arg {
	case "-?", "-h", "-help":
		return 1
	}
	if arg == "" || arg[0] != '-' {
		return 0
	}
	name := arg[1:]
	// If there's already "test.", drop it for now.
	name = strings.TrimPrefix(name, "test.")
	equals := strings.Index(name, "=")
	if equals >= 0 {
		name = name[:equals]
	}
	for _, f := range testFlagDefn {
		if name == f.name {
			// Booleans are special because they have modes -x, -x=true, -x=false.
			if !f.isBool && equals < 0 {
				return 2
			}
			return 1
		}
	}
	return 0
}
