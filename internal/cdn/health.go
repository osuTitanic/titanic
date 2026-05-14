package cdn

import (
	"context"
	"net/http"
)

func (c *Client) Health(ctx context.Context) (string, error) {
	return c.performRequestPlain(ctx, http.MethodGet, "/health", nil, nil, nil)
}
