package cdn

import (
	"context"
	"net/http"
)

type AdminSession struct {
	Name           string       `json:"name"`
	Prefixes       []string     `json:"prefixes"`
	Permissions    []Permission `json:"permissions"`
	UploadMode     string       `json:"upload_mode"`
	TrackDownloads bool         `json:"track_downloads"`
}

func (c *Client) AdminSession(ctx context.Context) (*AdminSession, error) {
	var session AdminSession
	if err := c.performRequestJson(ctx, http.MethodGet, "/admin/session", nil, nil, nil, &session); err != nil {
		return nil, err
	}
	return &session, nil
}
