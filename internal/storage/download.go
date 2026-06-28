package storage

import (
	"io"
	"net/http"
	"time"
)

var defaultAssetUrls = map[string]string{
	// Default avatars
	"avatars/unknown": "https://github.com/osuTitanic/titanic/blob/main/.github/images/avatars/unknown.jpg?raw=true",
	"avatars/1":       "https://github.com/osuTitanic/titanic/blob/main/.github/images/avatars/banchobot.jpg?raw=true",
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

func downloadAssetStream(key string) (io.ReadCloser, error) {
	url, exists := defaultAssetUrls[key]
	if !exists {
		return nil, nil
	}
	return downloadStream(url)
}

func downloadStream(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "osu!")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, err
	}
	return resp.Body, nil
}
