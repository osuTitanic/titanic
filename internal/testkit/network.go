package testkit

import (
	"context"
	"net/http"
	"time"
)

const (
	internetCheckUrl     = "https://example.com"
	internetCheckTimeout = 3 * time.Second
)

// IsInternetAvailable reports whether an external http endpoint can be reached.
func IsInternetAvailable(ctx context.Context) bool {
	ctx, cancel := context.WithTimeout(ctx, internetCheckTimeout)
	defer cancel()

	return isInternetAvailable(ctx, http.DefaultClient, internetCheckUrl)
}

func isInternetAvailable(ctx context.Context, client *http.Client, endpoint string) bool {
	request, err := http.NewRequestWithContext(ctx, http.MethodHead, endpoint, nil)
	if err != nil {
		return false
	}

	response, err := client.Do(request)
	if err != nil {
		return false
	}
	response.Body.Close()
	return true
}
