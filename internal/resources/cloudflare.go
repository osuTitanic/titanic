package resources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/osuTitanic/titanic-go/internal/config"
)

// Cloudflare accepts at most 30 urls per request
const maxUrlsPerRequest = 30
const purgeEndpoint = "https://api.cloudflare.com/client/v4/zones/%s/purge_cache"

var cloudflareHttpClient = &http.Client{Timeout: 10 * time.Second}

func PurgeOsz(setId int, cfg *config.Config) error {
	if !cfg.CloudflarePurgeEnabled {
		return nil
	}
	affectedUrls := resolveOszUrls(setId, cfg)
	return PurgeUrls(affectedUrls, cfg)
}

func PurgeUrls(urls []string, cfg *config.Config) error {
	if !cfg.CloudflarePurgeEnabled {
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

	slog.With("component", "cloudflare").Debug("Purged cached urls", "count", len(filteredUrls))
	return nil
}

func purgeUrlChunk(urls []string, cfg *config.Config) error {
	payload := map[string]any{"files": urls}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf(purgeEndpoint, cfg.CloudflareZoneId),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", fmt.Sprintf("osuTitanic/titanic (%s)", cfg.DomainName))
	req.Header.Set("Authorization", "Bearer "+cfg.CloudflareApiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := cloudflareHttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to purge urls: %d", resp.StatusCode)
	}

	return nil
}

func resolveOszUrls(setId int, cfg *config.Config) []string {
	templates := []string{
		fmt.Sprintf("%s/d/{id}", cfg.OsuBaseUrl()),
		fmt.Sprintf("%s/d/{id}n", cfg.OsuBaseUrl()),
	}
	if len(cfg.CloudflarePurgeOszUrls) > 0 {
		templates = cfg.CloudflarePurgeOszUrls
	}

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
