// +build builtinassets

package ui

import (
	"embed"
	"net/http"
)

//go:embed static
var assets embed.FS

// Assets contains the project's assets.
var Assets = http.FS(assets)
