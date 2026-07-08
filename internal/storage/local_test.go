package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileStorageSetup(t *testing.T) {
	dataPath := t.TempDir()
	assetPath := filepath.Join(dataPath, "avatars")
	if err := os.MkdirAll(assetPath, 0755); err != nil {
		t.Fatalf("failed to create asset directory: %v", err)
	}

	// Add placeholder files for required assets, since we don't
	// want to make calls to github.com in unit tests.
	for key := range defaultAssetUrls {
		path := filepath.Join(dataPath, filepath.FromSlash(key))
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatalf("failed to create asset parent directory: %v", err)
		}
		if err := os.WriteFile(path, []byte("placeholder"), 0644); err != nil {
			t.Fatalf("failed to write placeholder asset: %v", err)
		}
	}

	storage := NewFileStorage(dataPath)
	if err := storage.Setup(); err != nil {
		t.Fatalf("Setup() failed: %v", err)
	}

	// Check if setup was successful
	for _, directory := range RequiredDirectories {
		path := filepath.Join(dataPath, directory)
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("expected storage directory %q to exist: %v", directory, err)
		}
		if !info.IsDir() {
			t.Fatalf("expected %q to be a directory", directory)
		}
	}

	// TODO: `Save`, `Exists`, `Get` & `Remove` roundtrip tests
}

// TODO: S3 storage backend test (requires mocking the S3 API)
