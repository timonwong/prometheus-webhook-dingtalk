package main

import (
	// The blank import is to make go modules happy.
	_ "github.com/go-bindata/go-bindata"
)

//go:generate go run github.com/go-bindata/go-bindata/go-bindata -mode 420 -modtime 1 -pkg deftmpl -o template/internal/deftmpl/bindata.go template/default.tmpl
