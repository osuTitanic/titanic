package discord

import (
	"context"
	"fmt"
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

func (o *Officer) CallContext(ctx context.Context, tag OfficerTag, content string, embeds ...Embed) error {
	webhookUrl, ok := o.WebhookUrl(tag)
	if !ok {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}

	var contentPointer *string
	if content != "" {
		contentPointer = stringPointer(content)
	}

	webhook := &Webhook{
		URL:     webhookUrl,
		Content: contentPointer,
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

func stringPointer(value string) *string {
	return &value
}
