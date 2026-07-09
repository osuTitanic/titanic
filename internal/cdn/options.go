package cdn

import (
	"net/http"
)

type Option func(*Client)

func WithAccessKey(accessKey string) Option {
	return func(c *Client) {
		c.accessKey = accessKey
	}
}

func WithHttpClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithUserAgent(userAgent string) Option {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}
