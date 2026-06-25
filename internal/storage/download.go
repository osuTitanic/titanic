package storage

import (
	"io"
	"net/http"
)

var defaultAssetUrls = map[string]string{
	// Default avatars
	"avatars/unknown": "https://github.com/osuTitanic/titanic/blob/main/.github/images/avatars/unknown.jpg?raw=true",
	"avatars/1":       "https://github.com/osuTitanic/titanic/blob/main/.github/images/avatars/banchobot.jpg?raw=true",
}

func downloadAssetStream(key string) (io.ReadCloser, error) {
	url, exists := defaultAssetUrls[key]
	if !exists {
		return nil, nil
	}
	return downloadStream(url)
}

func downloadStream(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
