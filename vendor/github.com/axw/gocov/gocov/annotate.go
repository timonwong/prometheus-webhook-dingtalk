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
	"flag"
	"fmt"
	"go/token"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/axw/gocov"
)

const (
	hitPrefix  = "    "
	missPrefix = "MISS"
	RED        = "\x1b[31;1m"
	GREEN      = "\x1b[32;1m"
	NONE       = "\x1b[0m"
)

var (
	annotateFlags       = flag.NewFlagSet("annotate", flag.ExitOnError)
	annotateCeilingFlag = annotateFlags.Float64(
		"ceiling", 101,
		"Annotate only functions whose coverage is less than the specified percentage")
	annotateColorFlag = annotateFlags.Bool(
		"color", false,
		"Differentiate coverage with color")
)

type packageList []*gocov.Package
type functionList []*gocov.Function

func (l packageList) Len() int {
	return len(l)
}

func (l packageList) Less(i, j int) bool {
	return l[i].Name < l[j].Name
}

func (l packageList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l functionList) Len() int {
	return len(l)
}

func (l functionList) Less(i, j int) bool {
	return l[i].Name < l[j].Name
}

func (l functionList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

type annotator struct {
	fset  *token.FileSet
	files map[string]*token.File
}

func percentReached(fn *gocov.Function) float64 {
	if len(fn.Statements) == 0 {
		return 0
	}
	var reached int
	for _, stmt := range fn.Statements {
		if stmt.Reached > 0 {
			reached++
		}
	}
	return float64(reached) / float64(len(fn.Statements)) * 100
}

func annotateSource() (rc int) {
	annotateFlags.Parse(os.Args[2:])
	if annotateFlags.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "missing coverage file\n")
		return 1
	}

	var data []byte
	var err error
	if filename := annotateFlags.Arg(0); filename == "-" {
		data, err = ioutil.ReadAll(os.Stdin)
	} else {
		data, err = ioutil.ReadFile(filename)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read coverage file: %s\n", err)
		return 1
	}

	packages, err := unmarshalJson(data)
	if err != nil {
		fmt.Fprintf(
			os.Stderr, "failed to unmarshal coverage data: %s\n", err)
		return 1
	}

	// Sort packages, functions by name.
	sort.Sort(packageList(packages))
	for _, pkg := range packages {
		sort.Sort(functionList(pkg.Functions))
	}

	a := &annotator{}
	a.fset = token.NewFileSet()
	a.files = make(map[string]*token.File)

	var regexps []*regexp.Regexp
	for _, arg := range annotateFlags.Args()[1:] {
		re, err := regexp.Compile(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to compile %q as a regular expression, ignoring\n", arg)
		} else {
			regexps = append(regexps, re)
		}
	}
	if len(regexps) == 0 {
		regexps = append(regexps, regexp.MustCompile("."))
	}
	for _, pkg := range packages {
		for _, fn := range pkg.Functions {
			if percentReached(fn) >= *annotateCeilingFlag {
				continue
			}
			name := pkg.Name + "/" + fn.Name
			for _, regexp := range regexps {
				if regexp.FindStringIndex(name) != nil {
					err := a.printFunctionSource(fn)
					if err != nil {
						fmt.Fprintf(os.Stderr, "warning: failed to annotate function %q\n", name)
					}
					break
				}
			}
		}
	}
	return
}

func (a *annotator) printFunctionSource(fn *gocov.Function) error {
	// Load the file for line information. Probably overkill, maybe
	// just compute the lines from offsets in here.
	setContent := false
	file := a.files[fn.File]
	if file == nil {
		info, err := os.Stat(fn.File)
		if err != nil {
			return err
		}
		file = a.fset.AddFile(fn.File, a.fset.Base(), int(info.Size()))
		setContent = true
	}

	data, err := ioutil.ReadFile(fn.File)
	if err != nil {
		return err
	}
	if setContent {
		// This processes the content and records line number info.
		file.SetLinesForContent(data)
	}

	statements := fn.Statements[:]
	lineno := file.Line(file.Pos(fn.Start))
	lines := strings.Split(string(data)[fn.Start:fn.End], "\n")
	linenoWidth := int(math.Log10(float64(lineno+len(lines)))) + 1
	fmt.Println()
	for i, line := range lines {
		// Go through statements one at a time, seeing if we've hit
		// them or not.
		//
		// The prefix approach isn't perfect, as it doesn't
		// distinguish multiple statements per line. It'll have to
		// do for now. We could do fancy ANSI colouring later.
		lineno := lineno + i
		statementFound := false
		hit := false
		for j := 0; j < len(statements); j++ {
			start := file.Line(file.Pos(statements[j].Start))
			// FIXME instrumentation no longer records statements
			// in line order, as function literals are processed
			// after the body of a function. If/when that's changed,
			// we can go back to checking just the first statement
			// in each loop.
			if start == lineno {
				statementFound = true
				if !hit && statements[j].Reached > 0 {
					hit = true
				}
				statements = append(statements[:j], statements[j+1:]...)
			}
		}
		if *annotateColorFlag {
			color := NONE
			if statementFound && !hit {
				color = RED
			}
			fmt.Printf("%s%*d \t%s%s\n", color, linenoWidth, lineno, line, NONE)
		} else {
			hitmiss := hitPrefix
			if statementFound && !hit {
				hitmiss = missPrefix
			}
			fmt.Printf("%*d %s\t%s\n", linenoWidth, lineno, hitmiss, line)
		}
	}
	fmt.Println()

	return nil
}
