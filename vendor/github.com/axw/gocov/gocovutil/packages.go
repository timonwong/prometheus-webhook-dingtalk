package gocovutil

import (
	"encoding/json"
	"github.com/axw/gocov"
	"io/ioutil"
	"os"
	"sort"
)

// Packages represents a set of gocov.Package structures.
// The "AddPackage" method may be used to merge package
// coverage results into the set.
type Packages []*gocov.Package

// AddPackage adds a package's coverage information to the
func (ps *Packages) AddPackage(p *gocov.Package) {
	i := sort.Search(len(*ps), func(i int) bool {
		return (*ps)[i].Name >= p.Name
	})
	if i < len(*ps) && (*ps)[i].Name == p.Name {
		(*ps)[i].Accumulate(p)
	} else {
		head := (*ps)[:i]
		tail := append([]*gocov.Package{p}, (*ps)[i:]...)
		*ps = append(head, tail...)
	}
}

// ReadPackages takes a list of filenames and parses their
// contents as a Packages object.
//
// The special filename "-" may be used to indicate standard input.
// Duplicate filenames are ignored.
func ReadPackages(filenames []string) (ps Packages, err error) {
	copy_ := make([]string, len(filenames))
	copy(copy_, filenames)
	filenames = copy_
	sort.Strings(filenames)

	// Eliminate duplicates.
	unique := []string{filenames[0]}
	if len(filenames) > 1 {
		for _, f := range filenames[1:] {
			if f != unique[len(unique)-1] {
				unique = append(unique, f)
			}
		}
	}

	// Open files.
	var files []*os.File
	for _, f := range filenames {
		if f == "-" {
			files = append(files, os.Stdin)
		} else {
			file, err := os.Open(f)
			if err != nil {
				return nil, err
			}
			defer file.Close()
			files = append(files, os.Stdin)
		}
	}

	// Parse the files, accumulate Packages.
	for _, file := range files {
		data, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		result := &struct{ Packages []*gocov.Package }{}
		err = json.Unmarshal(data, result)
		if err != nil {
			return nil, err
		}
		for _, p := range result.Packages {
			ps.AddPackage(p)
		}
	}
	return ps, nil
}
