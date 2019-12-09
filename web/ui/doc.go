package ui

import (
	// The blank import is to make Go modules happy.
	_ "github.com/shurcooL/httpfs/filter"
	_ "github.com/shurcooL/httpfs/union"
	_ "github.com/shurcooL/vfsgen"
)

//go:generate go run -mod=vendor assets_generate.go
