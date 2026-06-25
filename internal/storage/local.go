package storage

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type FileStorage struct {
	dataPath string
}

func NewFileStorage(dataPath string) Storage {
	return &FileStorage{dataPath: dataPath}
}

func (storage *FileStorage) Setup() error {
	err := storage.CreateDefaultFolders()
	if err != nil {
		return err
	}

	err = storage.DownloadDefaultAssets()
	if err != nil {
		return err
	}
	return nil
}

func (storage *FileStorage) Read(key string, folder string) ([]byte, error) {
	path := fmt.Sprintf("%s/%s/%s", storage.dataPath, folder, key)
	return os.ReadFile(path)
}

func (storage *FileStorage) ReadStream(key string, folder string) (io.ReadSeekCloser, error) {
	path := fmt.Sprintf("%s/%s/%s", storage.dataPath, folder, key)
	return os.Open(path)
}

func (storage *FileStorage) Save(key string, folder string, data []byte) error {
	path := fmt.Sprintf("%s/%s", storage.dataPath, folder)
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}

	os.WriteFile(fmt.Sprintf("%s/%s", path, key), data, os.ModePerm)
	return nil
}

func (storage *FileStorage) SaveStream(key string, folder string, stream io.Reader) error {
	filePath := fmt.Sprintf("%s/%s/%s", storage.dataPath, folder, key)
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	buffer := make([]byte, 256*1024)
	_, err = io.CopyBuffer(file, stream, buffer)
	return err
}

func (storage *FileStorage) SaveUrl(key string, directory string, url string) error {
	stream, err := downloadStream(url)
	if err != nil {
		return fmt.Errorf("failed to download url content %q: %w", url, err)
	}
	if stream == nil {
		return fmt.Errorf("no download stream found for url %q", url)
	}
	defer stream.Close()

	err = storage.SaveStream(key, directory, stream)
	if err != nil {
		return err
	}
	return nil
}

func (storage *FileStorage) Exists(key string, folder string) bool {
	path := fmt.Sprintf("%s/%s/%s", storage.dataPath, folder, key)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func (storage *FileStorage) Remove(key string, folder string) error {
	path := fmt.Sprintf("%s/%s/%s", storage.dataPath, folder, key)
	return os.Remove(path)
}

// CreateDefaultFolders ensures that all required storage folders exist
func (storage *FileStorage) CreateDefaultFolders() error {
	for _, directory := range RequiredDirectories {
		folder := fmt.Sprintf(
			"%s/%s",
			storage.dataPath, directory,
		)

		if _, err := os.Stat(folder); !os.IsNotExist(err) {
			slog.Debug("Storage directory already exists, skipping creation", slog.String("path", folder))
			continue
		}

		err := os.MkdirAll(folder, 0755)
		if err != nil {
			return fmt.Errorf("failed to create storage directory %s: %w", folder, err)
		}

		slog.Info("Created storage directory", slog.String("path", folder))
	}
	return nil
}

// DownloadDefaultAssets downloads all default assets if they do not already exist
func (storage *FileStorage) DownloadDefaultAssets() error {
	for assetUrl := range defaultAssetUrls {
		parts := strings.SplitN(assetUrl, "/", 2)
		if len(parts) != 2 {
			slog.Warn("Invalid asset URL, skipping download", slog.String("url", assetUrl))
			continue
		}
		folder := parts[0]
		key := parts[1]

		// Check if asset already exists
		if storage.Exists(key, folder) {
			slog.Debug("Asset already exists, skipping download", slog.String("path", assetUrl))
			continue
		}

		stream, err := downloadAssetStream(assetUrl)
		if err != nil {
			return fmt.Errorf("failed to get download stream for %s: %w", assetUrl, err)
		}
		if stream == nil {
			return fmt.Errorf("no download stream found for %s", assetUrl)
		}
		defer stream.Close()

		err = storage.SaveStream(key, folder, stream)
		if err != nil {
			return fmt.Errorf("failed to save download %s to storage: %w", assetUrl, err)
		}
		slog.Info("Downloaded asset to storage", slog.String("path", assetUrl))
	}
	return nil
}
