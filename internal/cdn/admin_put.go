package cdn

import (
	"context"
	"errors"
	"io"
	"net/http"
)

type UploadResponse struct {
	Key  string `json:"key"`
	ETag string `json:"etag"`
}

type UploadFileOptions struct {
	ContentType        string
	CacheControl       string
	ContentDisposition string
}

func (c *Client) AdminUploadFile(ctx context.Context, objectKey string, body io.Reader, options *UploadFileOptions) (*UploadResponse, error) {
	if body == nil {
		return nil, errors.New("cdn: upload body is required")
	}

	objectPath, err := objectPath(objectKey)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	if options != nil {
		setHeader(headers, "Content-Type", options.ContentType)
		setHeader(headers, "Cache-Control", options.CacheControl)
		setHeader(headers, "Content-Disposition", options.ContentDisposition)
	}

	var upload UploadResponse
	if err := c.performRequestJson(ctx, http.MethodPut, "/admin/files/"+objectPath, nil, headers, body, &upload); err != nil {
		return nil, err
	}
	return &upload, nil
}
