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

package cov

import (
	"fmt"
	"github.com/axw/gocov"
	"go/token"
	"html"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const (
	hitPrefix  = "    "
	missPrefix = "MISS"
)

type annotator struct {
	fset  *token.FileSet
	files map[string]*token.File
}

func annotateFunctionSource(w io.Writer, fn *gocov.Function) {
	if fn == nil {
		panic("nil function to annotate")
	}
	a := &annotator{}
	a.fset = token.NewFileSet()
	a.files = make(map[string]*token.File)
	a.printFunctionSource(w, fn)
}

func (a *annotator) printFunctionSource(w io.Writer, fn *gocov.Function) error {
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
	fmt.Fprintln(w)
	fmt.Fprintf(w, "<div class=\"funcname\" id=\"fn_%s\">func %s</div>", fn.Name, fn.Name)
	fmt.Fprintf(w, "<div class=\"info\"><a href=\"#s_fn_%s\">Back</a><p>In <code>%s</code>:</p></div>",
		fn.Name, fn.File)
	fmt.Fprintln(w, "<table class=\"listing\">")
	for i, line := range lines {
		lineno := lineno + i
		statementFound := false
		hit := false
		for j := 0; j < len(statements); j++ {
			start := file.Line(file.Pos(statements[j].Start))
			if start == lineno {
				statementFound = true
				if !hit && statements[j].Reached > 0 {
					hit = true
				}
				statements = append(statements[:j], statements[j+1:]...)
			}
		}
		hitmiss := hitPrefix
		if statementFound && !hit {
			hitmiss = missPrefix
		}
		tr := "<tr"
		if hitmiss == missPrefix {
			tr += ` class="miss">`
		} else {
			tr += ">"
		}
		fmt.Fprintf(w, "%s<td>%d</td><td><code><pre>%s</pre></code></td></tr>", tr, lineno,
			html.EscapeString(strings.Replace(line, "\t", "        ", -1)))
	}
	fmt.Fprintln(w, "</table>")

	return nil
}
