//go:build builtinassets
// +build builtinassets

package template

import (
	"embed"
	"net/http"

	"github.com/shurcooL/httpfs/union"
)

//go:embed *.tmpl
var assets embed.FS

// Assets contains the project's assets.
var Assets = func() http.FileSystem {
	return union.New(map[string]http.FileSystem{
		"/templates": http.FS(assets),
	})
}()
