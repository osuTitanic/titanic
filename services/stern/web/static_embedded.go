//go:build static_embedded

// NOTE: By default the "static" folder is not being embedded into the binary.
//       If you want to embed the static assets, build with the "static_embedded" tag:
//       go build -tags static_embedded

package web

import (
	"embed"
	"io/fs"
	"path"
	"strings"
)

//go:embed static
var Static embed.FS

func StaticFS(assetPath string) (fs.FS, error) {
	cleaned := strings.Trim(assetPath, "/")
	if cleaned == "" {
		return fs.Sub(Static, "static")
	}

	target := path.Join("static", cleaned)
	info, err := fs.Stat(Static, target)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return fs.Sub(Static, target)
	}

	return fs.Sub(Static, "static")
}
