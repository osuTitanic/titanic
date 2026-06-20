package resources

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

// TODO: we could probably refactor this to be the "global standard"
// 	     http implementation for the entire codebase

// retryStatusCodes are the status codes that should be retried by our http client.
var retryStatusCodes = map[int]struct{}{
	http.StatusInternalServerError: {}, // 500
	http.StatusBadGateway:          {}, // 502
	http.StatusServiceUnavailable:  {}, // 503
	http.StatusGatewayTimeout:      {}, // 504
}

// httpSession is a thin wrapper around http.Client that
// attaches our user agent to every outgoing request.
type httpSession struct {
	client    *http.Client
	userAgent string
}

func (session *httpSession) Get(url string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", session.userAgent)
	return session.client.Do(request)
}

// createHttpSession creates an http client that retries on server
// errors & identifies itself with the given user agent.
func createHttpSession(userAgent string) *httpSession {
	client := retryablehttp.NewClient()
	client.RetryMax = 4
	client.RetryWaitMin = 300 * time.Millisecond
	client.CheckRetry = retryPolicy
	client.Logger = nil
	return &httpSession{client: client.StandardClient(), userAgent: userAgent}
}

// retryPolicy defines the retry policy for our http client.
// It retries on server errors and respects context cancellation.
func retryPolicy(ctx context.Context, response *http.Response, err error) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	if err != nil {
		return true, nil
	}

	_, shouldRetry := retryStatusCodes[response.StatusCode]
	return shouldRetry, nil
}
