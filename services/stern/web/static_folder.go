//go:build !static_embedded

// NOTE: By default the "static" folder is not being embedded into the binary.
//       If you want to embed the static assets, build with the "static_embedded" tag:
//       go build -tags static_embedded

package web

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const StaticDirectoryEnv = "STERN_STATIC_DIR"

func StaticFS(assetPath string) (fs.FS, error) {
	root := StaticRootPath()
	cleaned := strings.Trim(assetPath, "/")
	if cleaned == "" {
		return os.DirFS(root), nil
	}

	target := filepath.Join(root, filepath.FromSlash(cleaned))
	info, err := os.Stat(target)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return os.DirFS(target), nil
	}

	return os.DirFS(root), nil
}

func StaticRootPath() string {
	if root := strings.TrimSpace(os.Getenv(StaticDirectoryEnv)); root != "" {
		return root
	}

	// Assume we are running from the root folder by default
	// I'm too lazy to develop another dynamic folder resolve logic thingy
	return filepath.Join("services", "stern", "web", "static")
}
