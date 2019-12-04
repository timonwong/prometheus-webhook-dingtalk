// +build tools

// Package tools manages development tool versions through the module system.
//
// See https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
package tools

import (
	_ "github.com/axw/gocov/gocov"
	_ "github.com/go-bindata/go-bindata/go-bindata"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/matm/gocov-html"
	_ "github.com/prometheus/promu"
	_ "golang.org/x/tools/cmd/goimports"
)
