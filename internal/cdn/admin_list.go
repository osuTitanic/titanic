package cdn

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type ListFilesOptions struct {
	Prefix string
	Limit  int
	Cursor string
}

type ListResponse struct {
	Prefix     string     `json:"prefix"`
	Items      []ListItem `json:"items"`
	NextCursor string     `json:"next_cursor"`
}

type ListItem struct {
	Key           string    `json:"key"`
	Size          int64     `json:"size"`
	ETag          string    `json:"etag"`
	LastModified  time.Time `json:"last_modified"`
	DownloadCount *uint64   `json:"download_count,omitempty"`
}

func (c *Client) AdminListFiles(ctx context.Context, params ListFilesOptions) (*ListResponse, error) {
	query := url.Values{}
	query.Set("prefix", params.Prefix)
	if params.Limit > 0 {
		query.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Cursor != "" {
		query.Set("cursor", params.Cursor)
	}

	var list ListResponse
	if err := c.performRequestJson(ctx, http.MethodGet, "/admin/files", query, nil, nil, &list); err != nil {
		return nil, err
	}
	return &list, nil
}
