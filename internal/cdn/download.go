package cdn

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type ObjectResponse struct {
	StatusCode int
	Header     http.Header
	Body       io.ReadCloser
}

func newObjectResponse(resp *http.Response) *ObjectResponse {
	return &ObjectResponse{
		StatusCode: resp.StatusCode,
		Header:     resp.Header.Clone(),
		Body:       resp.Body,
	}
}

func (r *ObjectResponse) Location() string {
	return r.Header.Get("Location")
}

func (r *ObjectResponse) ContentType() string {
	return r.Header.Get("Content-Type")
}

func (r *ObjectResponse) ContentRange() string {
	return r.Header.Get("Content-Range")
}

func (r *ObjectResponse) ContentLength() int64 {
	value := r.Header.Get("Content-Length")
	if value == "" {
		return -1
	}
	length, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return -1
	}
	return length
}

func (r *ObjectResponse) AcceptRanges() string {
	return r.Header.Get("Accept-Ranges")
}

func (r *ObjectResponse) ETag() string {
	return r.Header.Get("ETag")
}

func (r *ObjectResponse) LastModified() string {
	return r.Header.Get("Last-Modified")
}

func (r *ObjectResponse) CacheControl() string {
	return r.Header.Get("Cache-Control")
}

type ObjectRequestOptions struct {
	Range string
}

func (c *Client) GetObject(ctx context.Context, objectKey string, options *ObjectRequestOptions) (*ObjectResponse, error) {
	objectPath, err := objectPath(objectKey)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	if options != nil {
		setHeader(headers, "Range", options.Range)
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/"+objectPath, nil, headers, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cdn: get object: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusPartialContent:
		objectResponse := newObjectResponse(resp)
		return objectResponse, nil
	case http.StatusTemporaryRedirect:
		objectResponse := newObjectResponse(resp)
		resp.Body.Close()
		objectResponse.Body = nil
		return objectResponse, nil
	default:
		defer resp.Body.Close()
		return nil, c.resolveResponseError(req, resp)
	}
}
