package discord

import (
	"bytes"
	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/osuTitanic/titanic/internal/config"
)

type Officer struct {
	Url     string
	Timeout time.Duration
}

type OfficerTag string

const (
	OfficerTagGeneral         OfficerTag = "general"
	OfficerTagRegistration    OfficerTag = "registration"
	OfficerTagScoreSubmission OfficerTag = "score-submission"
	OfficerTagReports         OfficerTag = "reports"
)

// TODO: Add separate webhook URLs for each tag

func NewOfficer(webhookUrl string) *Officer {
	return &Officer{
		Url:     strings.TrimSpace(webhookUrl),
		Timeout: 5 * time.Second,
	}
}

func NewOfficerFromConfig(cfg *config.Config) *Officer {
	if cfg == nil {
		return NewOfficer("")
	}
	return NewOfficer(cfg.OfficerWebhookUrl)
}

func (o *Officer) Enabled(tag OfficerTag) bool {
	_, ok := o.WebhookUrl(tag)
	return ok
}

func (o *Officer) WebhookUrl(tag OfficerTag) (string, bool) {
	if o == nil || o.Url == "" {
		return "", false
	}
	return o.Url, true
}

func (o *Officer) Call(tag OfficerTag, content string, embeds ...Embed) error {
	return o.CallContext(context.Background(), tag, content, embeds...)
}

func (o *Officer) CallWithFile(tag OfficerTag, content, filename string, data []byte, embeds ...Embed) error {
	return o.CallWithFileContext(context.Background(), tag, content, filename, data, embeds...)
}

func (o *Officer) CallContext(ctx context.Context, tag OfficerTag, content string, embeds ...Embed) error {
	return o.callContext(ctx, tag, content, nil, embeds...)
}

func (o *Officer) CallWithFileContext(ctx context.Context, tag OfficerTag, content, filename string, data []byte, embeds ...Embed) error {
	file := &File{Name: filename, Content: bytes.NewReader(data)}
	return o.callContext(ctx, tag, content, file, embeds...)
}

func (o *Officer) CallPanic(recovered any, context ...Field) error {
	embed := Embed{Title: new("Unhandled Panic"), Color: new(0xFF0000)}
	embed.AddField("Error", fmt.Sprint(recovered), false)

	for _, field := range context {
		embed.Fields = append(embed.Fields, field)
	}

	return o.CallWithFile(
		OfficerTagGeneral,
		"",
		"stacktrace.txt",
		debug.Stack(),
		embed,
	)
}

func (o *Officer) callContext(ctx context.Context, tag OfficerTag, content string, file *File, embeds ...Embed) error {
	webhookUrl, ok := o.WebhookUrl(tag)
	if !ok {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}

	var contentPointer *string
	if content != "" {
		contentPointer = new(content)
	}

	webhook := &Webhook{
		URL:     webhookUrl,
		Content: contentPointer,
		File:    file,
		Embeds:  embeds,
	}

	if o.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, o.Timeout)
		defer cancel()
	}

	if err := webhook.PostContext(ctx); err != nil {
		return fmt.Errorf("officer %q webhook failed: %w", normalizeOfficerTag(tag), err)
	}
	return nil
}

func normalizeOfficerTag(tag OfficerTag) OfficerTag {
	value := strings.ToLower(strings.TrimSpace(string(tag)))
	if value == "" {
		return OfficerTagGeneral
	}
	return OfficerTag(value)
}
