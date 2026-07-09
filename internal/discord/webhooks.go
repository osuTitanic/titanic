package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

const WebhookUserAgent = "osuTitanic/titanic"

type Footer struct {
	Text         string  `json:"text"`
	IconURL      *string `json:"icon_url,omitempty"`
	ProxyIconURL *string `json:"proxy_icon_url,omitempty"`
}

type Image struct {
	URL      string  `json:"url"`
	ProxyURL *string `json:"proxy_url,omitempty"`
	Height   *int    `json:"height,omitempty"`
	Width    *int    `json:"width,omitempty"`
}

type Thumbnail struct {
	URL      string  `json:"url"`
	ProxyURL *string `json:"proxy_url,omitempty"`
	Height   *int    `json:"height,omitempty"`
	Width    *int    `json:"width,omitempty"`
}

type Video struct {
	URL    string `json:"url"`
	Height *int   `json:"height,omitempty"`
	Width  *int   `json:"width,omitempty"`
}

type Provider struct {
	URL  *string `json:"url,omitempty"`
	Name *string `json:"name,omitempty"`
}

type Author struct {
	Name         *string `json:"name,omitempty"`
	URL          *string `json:"url,omitempty"`
	IconURL      *string `json:"icon_url,omitempty"`
	ProxyIconURL *string `json:"proxy_icon_url,omitempty"`
}

type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

type Embed struct {
	Title       *string    `json:"title,omitempty"`
	Type        *string    `json:"type,omitempty"`
	Description *string    `json:"description,omitempty"`
	URL         *string    `json:"url,omitempty"`
	Timestamp   *time.Time `json:"timestamp,omitempty"`
	Color       *int       `json:"color,omitempty"`
	Footer      *Footer    `json:"footer,omitempty"`
	Image       *Image     `json:"image,omitempty"`
	Thumbnail   *Thumbnail `json:"thumbnail,omitempty"`
	Video       *Video     `json:"video,omitempty"`
	Provider    *Provider  `json:"provider,omitempty"`
	Author      *Author    `json:"author,omitempty"`
	Fields      []Field    `json:"fields,omitempty"`
}

func (e *Embed) AddField(name, value string, inline bool) {
	e.Fields = append(e.Fields, Field{
		Name:   name,
		Value:  value,
		Inline: inline,
	})
}

type File struct {
	Name    string
	Content io.Reader
}

type Webhook struct {
	URL       string
	Content   *string
	Username  *string
	AvatarURL *string
	TTS       *bool
	File      *File
	Embeds    []Embed
}

func (w *Webhook) AddEmbed(embed Embed) {
	w.Embeds = append(w.Embeds, embed)
}

func (w *Webhook) SetFile(name string, data []byte) {
	w.File = &File{Name: name, Content: bytes.NewReader(data)}
}

func (w *Webhook) SetFileReader(name string, content io.Reader) {
	w.File = &File{Name: name, Content: content}
}

func (w *Webhook) Post() error {
	return w.PostContext(context.Background())
}

func (w *Webhook) PostContext(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	p, err := w.buildPayload()
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	if w.File == nil {
		return w.postJson(ctx, jsonData)
	}
	return w.postMultipart(ctx, jsonData)
}

type payload struct {
	Content   *string `json:"content,omitempty"`
	Username  *string `json:"username,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	TTS       *bool   `json:"tts,omitempty"`
	Embeds    []Embed `json:"embeds"`
}

func (w *Webhook) buildPayload() (*payload, error) {
	// Ensure we have a valid payload
	hasContent := w.Content != nil && *w.Content != ""
	hasFile := w.File != nil && w.File.Content != nil
	hasEmbeds := len(w.Embeds) > 0

	if !hasContent && !hasFile && !hasEmbeds {
		return nil, fmt.Errorf(
			"webhook must contain at least content, a file, or an embed",
		)
	}

	// TODO: Truncate content to 2k characters
	return &payload{
		Content:   w.Content,
		Username:  w.Username,
		AvatarURL: w.AvatarURL,
		TTS:       w.TTS,
		Embeds:    w.Embeds,
	}, nil
}

func (w *Webhook) postJson(ctx context.Context, jsonData []byte) error {
	req, err := http.NewRequestWithContext(ctx, "POST", w.URL, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", WebhookUserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// TODO: Read response data
		return fmt.Errorf(
			"webhook returned status %d",
			resp.StatusCode,
		)
	}
	return nil
}

func (w *Webhook) postMultipart(ctx context.Context, jsonData []byte) error {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	if err := writer.WriteField("payload_json", string(jsonData)); err != nil {
		return fmt.Errorf("failed to write payload field: %w", err)
	}

	part, err := writer.CreateFormFile("file", w.File.Name)
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, w.File.Content); err != nil {
		return fmt.Errorf("failed to write file data: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", w.URL, &buf)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Agent", WebhookUserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// TODO: Read response data
		return fmt.Errorf(
			"webhook returned status %d",
			resp.StatusCode,
		)
	}
	return nil
}
