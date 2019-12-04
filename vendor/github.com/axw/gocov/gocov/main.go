// Copyright (c) 2012 The Gocov Authors.
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
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/axw/gocov"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n\n\tgocov command [arguments]\n\n")
	fmt.Fprintf(os.Stderr, "The commands are:\n\n")
	fmt.Fprintf(os.Stderr, "\tannotate\n")
	fmt.Fprintf(os.Stderr, "\tconvert\n")
	fmt.Fprintf(os.Stderr, "\treport\n")
	fmt.Fprintf(os.Stderr, "\ttest\n")
	fmt.Fprintf(os.Stderr, "\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func marshalJson(packages []*gocov.Package) ([]byte, error) {
	return json.Marshal(struct{ Packages []*gocov.Package }{packages})
}

func unmarshalJson(data []byte) (packages []*gocov.Package, err error) {
	result := &struct{ Packages []*gocov.Package }{}
	err = json.Unmarshal(data, result)
	if err == nil {
		packages = result.Packages
	}
	return
}

func main() {
	flag.Usage = usage
	flag.Parse()

	command := ""
	if flag.NArg() > 0 {
		command = flag.Arg(0)
		switch command {
		case "convert":
			if flag.NArg() <= 1 {
				fmt.Fprintln(os.Stderr, "missing cover profile")
				os.Exit(1)
			}
			if err := convertProfiles(flag.Args()[1:]...); err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				os.Exit(1)
			}
		case "annotate":
			os.Exit(annotateSource())
		case "report":
			os.Exit(reportCoverage())
		case "test":
			if err := runTests(flag.Args()[1:]); err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				os.Exit(1)
			}
		default:
			fmt.Fprintf(os.Stderr, "Unknown command: %#q\n\n", command)
			usage()
		}
	} else {
		usage()
	}
}
