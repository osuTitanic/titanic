package cdn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	baseUrl    *url.URL
	httpClient *http.Client
	accessKey  string
	userAgent  string
}

func (c *Client) BaseUrl() string {
	return c.baseUrl.String()
}

func NewDefaultClient(options ...Option) (*Client, error) {
	return NewClient(DefaultStreamingBaseUrl, options...)
}

func NewClient(baseUrl string, options ...Option) (*Client, error) {
	parsedUrl, err := url.Parse(baseUrl)
	if err != nil {
		return nil, fmt.Errorf("cdn: parse base url: %w", err)
	}
	if parsedUrl.Scheme == "" || parsedUrl.Host == "" {
		return nil, errors.New("cdn: base url must include scheme and host")
	}

	client := &Client{
		baseUrl:    parsedUrl,
		userAgent:  DefaultUserAgent,
		httpClient: defaultHttpClient(),
	}
	for _, option := range options {
		option(client)
	}
	if client.httpClient == nil {
		return nil, errors.New("cdn: http client is nil")
	}
	return client, nil
}

// i really hate these long ass methods

func (c *Client) newRequest(
	ctx context.Context,
	method string,
	path string,
	query url.Values,
	headers http.Header,
	body io.Reader,
) (*http.Request, error) {
	requestUrl := *c.baseUrl
	requestUrl.Path = joinUrlPath(c.baseUrl.Path, path)
	requestUrl.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, method, requestUrl.String(), body)
	if err != nil {
		return nil, fmt.Errorf("cdn: create request: %w", err)
	}

	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	if c.accessKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessKey)
	}
	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	return req, nil
}

func (c *Client) performRequestJson(
	ctx context.Context,
	method string,
	path string,
	query url.Values,
	headers http.Header,
	body io.Reader,
	out any,
) error {
	req, err := c.newRequest(ctx, method, path, query, headers, body)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("cdn: %s %s: %w", method, path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return c.resolveResponseError(req, resp)
	}
	if out == nil || resp.StatusCode == http.StatusNoContent {
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("cdn: decode %s %s response: %w", method, path, err)
	}
	return nil
}

func (c *Client) performRequestPlain(
	ctx context.Context,
	method string,
	path string,
	query url.Values,
	headers http.Header,
	body io.Reader,
) (string, error) {
	req, err := c.newRequest(ctx, method, path, query, headers, body)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("cdn: %s %s: %w", method, path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", c.resolveResponseError(req, resp)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cdn: read %s %s response: %w", method, path, err)
	}
	return string(responseBody), nil
}

func (c *Client) resolveResponseError(req *http.Request, resp *http.Response) error {
	// Ensure we don't read excessively large error responses
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("cdn: read error response: %w", err)
	}

	apiError := &ApiError{
		Method:     req.Method,
		Path:       req.URL.RequestURI(),
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Body:       body,
	}

	if !isResponseJson(resp) {
		apiError.Message = strings.TrimSpace(string(body))
		return apiError
	}

	var errorResponse ErrorResponse

	// Try to parse json error response
	if err := json.Unmarshal(body, &errorResponse); err == nil {
		apiError.ErrorCode = errorResponse.Error
		apiError.Message = errorResponse.Message
		return apiError
	}

	apiError.Message = strings.TrimSpace(string(body))
	return apiError
}
