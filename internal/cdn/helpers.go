package cdn

import (
	"errors"
	"net/http"
	"slices"
	"strings"
	"time"
)

func setHeader(headers http.Header, key string, value string) {
	if value != "" {
		headers.Set(key, value)
	}
}

func isHeaderJson(contentType string) bool {
	return strings.HasPrefix(strings.ToLower(contentType), "application/json")
}

func isResponseJson(resp *http.Response) bool {
	return isHeaderJson(resp.Header.Get("Content-Type"))
}

func defaultHttpClient() *http.Client {
	return &http.Client{Timeout: 30 * time.Second}
}

func joinUrlPath(basePath string, requestPath string) string {
	basePath = strings.TrimRight(basePath, "/")
	if requestPath == "" {
		return basePath
	}
	return basePath + "/" + strings.TrimLeft(requestPath, "/")
}

func objectPath(objectKey string) (string, error) {
	if objectKey == "" {
		return "", errors.New("cdn: object key is required")
	}
	if strings.HasPrefix(objectKey, "/") {
		return "", errors.New("cdn: object key must be relative")
	}

	segments := strings.Split(objectKey, "/")
	if slices.Contains(segments, "") {
		return "", errors.New("cdn: object key contains an empty path segment")
	}
	return objectKey, nil
}
