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

// Package gocov is a code coverage analysis tool for Go.
package gocov

import (
	"fmt"
)

type Package struct {
	// Name is the canonical path of the package.
	Name string

	// Functions is a list of functions registered with this package.
	Functions []*Function
}

type Function struct {
	// Name is the name of the function. If the function has a receiver, the
	// name will be of the form T.N, where T is the type and N is the name.
	Name string

	// File is the full path to the file in which the function is defined.
	File string

	// Start is the start offset of the function's signature.
	Start int

	// End is the end offset of the function.
	End int

	// statements registered with this function.
	Statements []*Statement
}

type Statement struct {
	// Start is the start offset of the statement.
	Start int

	// End is the end offset of the statement.
	End int

	// Reached is the number of times the statement was reached.
	Reached int64
}

// Accumulate will accumulate the coverage information from the provided
// Package into this Package.
func (p *Package) Accumulate(p2 *Package) error {
	if p.Name != p2.Name {
		return fmt.Errorf("Names do not match: %q != %q", p.Name, p2.Name)
	}
	if len(p.Functions) != len(p2.Functions) {
		return fmt.Errorf("Function counts do not match: %d != %d", len(p.Functions), len(p2.Functions))
	}
	for i, f := range p.Functions {
		err := f.Accumulate(p2.Functions[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// Accumulate will accumulate the coverage information from the provided
// Function into this Function.
func (f *Function) Accumulate(f2 *Function) error {
	if f.Name != f2.Name {
		return fmt.Errorf("Names do not match: %q != %q", f.Name, f2.Name)
	}
	if f.File != f2.File {
		return fmt.Errorf("Files do not match: %q != %q", f.File, f2.File)
	}
	if f.Start != f2.Start || f.End != f2.End {
		return fmt.Errorf("Source ranges do not match: %d-%d != %d-%d", f.Start, f.End, f2.Start, f2.End)
	}
	if len(f.Statements) != len(f2.Statements) {
		return fmt.Errorf("Number of statements do not match: %d != %d", len(f.Statements), len(f2.Statements))
	}
	for i, s := range f.Statements {
		err := s.Accumulate(f2.Statements[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// Accumulate will accumulate the coverage information from the provided
// Statement into this Statement.
func (s *Statement) Accumulate(s2 *Statement) error {
	if s.Start != s2.Start || s.End != s2.End {
		return fmt.Errorf("Source ranges do not match: %d-%d != %d-%d", s.Start, s.End, s2.Start, s2.End)
	}
	s.Reached += s2.Reached
	return nil
}
