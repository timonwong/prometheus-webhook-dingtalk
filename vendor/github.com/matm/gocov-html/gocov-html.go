// Copyright (c) 2013 Mathias Monnerville
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
	"flag"
	"github.com/matm/gocov-html/cov"
	"io"
	"log"
	"os"
)

func main() {
	var r io.Reader
	log.SetFlags(0)

	var s = flag.String("s", "", "path to custom CSS file")
	flag.Parse()

	switch flag.NArg() {
	case 0:
		r = os.Stdin
	case 1:
		var err error
		if r, err = os.Open(flag.Arg(0)); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Usage: %s data.json\n", os.Args[0])
	}

	if err := cov.HTMLReportCoverage(r, *s); err != nil {
		log.Fatal(err)
	}
}
