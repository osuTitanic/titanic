package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type CloudflarePurgeConfiguration struct {
	PurgeEnabled bool
	ZoneId       string
	ApiToken     string
	OszPurgeUrls []string
	UserAgent    string
}

func (cfg *CloudflarePurgeConfiguration) CanPurge() bool {
	return cfg.PurgeEnabled && cfg.ZoneId != "" && cfg.ApiToken != ""
}

// Cloudflare accepts at most 30 urls per request
const maxUrlsPerRequest = 30
const purgeEndpoint = "https://api.cloudflare.com/client/v4/zones/%s/purge_cache"

func PurgeOsz(setId int, cfg *CloudflarePurgeConfiguration) error {
	if !cfg.CanPurge() {
		return nil
	}
	affectedUrls := oszUrls(setId, cfg.OszPurgeUrls)
	return PurgeUrls(affectedUrls, cfg)
}

func PurgeUrls(urls []string, cfg *CloudflarePurgeConfiguration) error {
	if !cfg.CanPurge() {
		return nil
	}

	filteredUrls := filterUrls(urls)
	if len(filteredUrls) <= 0 {
		return nil
	}

	for _, chunk := range chunkedUrls(filteredUrls, maxUrlsPerRequest) {
		err := purgeUrlChunk(chunk, cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

func purgeUrlChunk(urls []string, cfg *CloudflarePurgeConfiguration) error {
	payload := map[string]any{"files": urls}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf(purgeEndpoint, cfg.ZoneId),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+cfg.ApiToken)
	req.Header.Set("Content-Type", "application/json")
	if cfg.UserAgent != "" {
		req.Header.Set("User-Agent", cfg.UserAgent)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to purge urls: %d", resp.StatusCode)
	}

	return nil
}

func oszUrls(setId int, templates []string) []string {
	urls := make([]string, 0, len(templates))
	for _, template := range templates {
		url := strings.ReplaceAll(template, "{id}", fmt.Sprintf("%d", setId))
		urls = append(urls, url)
	}
	return urls
}

func filterUrls(urls []string) []string {
	// Deduplicate & remove empty urls
	uniqueUrls := make(map[string]struct{})
	for _, url := range urls {
		if url != "" {
			uniqueUrls[url] = struct{}{}
		}
	}

	// Convert map keys back to slice
	filteredUrls := make([]string, 0, len(uniqueUrls))
	for url := range uniqueUrls {
		filteredUrls = append(filteredUrls, url)
	}
	return filteredUrls
}

func chunkedUrls(urls []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(urls); i += chunkSize {
		end := min(i+chunkSize, len(urls))
		chunks = append(chunks, urls[i:end])
	}
	return chunks
}
