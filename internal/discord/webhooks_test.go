package discord

import (
	"flag"
	"strings"
	"testing"
)

var testWebhookURL = flag.String("webhook-url", "", "Discord webhook URL used for integration tests")

func requireWebhookUrl(t *testing.T) string {
	t.Helper()

	if *testWebhookURL == "" {
		t.Skip("skipping integration test: provide -webhook-url")
	}
	return *testWebhookURL
}

func TestWebhookPostJson(t *testing.T) {
	url := requireWebhookUrl(t)
	content := "common-go webhook json test"

	w := &Webhook{
		URL:     url,
		Content: &content,
	}
	if err := w.Post(); err != nil {
		t.Fatalf("Post() error = %v", err)
	}
}

func TestWebhookPostMultipart(t *testing.T) {
	url := requireWebhookUrl(t)
	content := "common-go webhook multipart test"

	w := &Webhook{
		URL:     url,
		Content: &content,
	}
	w.SetFileReader("test.txt", strings.NewReader("webhook multipart test payload"))

	if err := w.Post(); err != nil {
		t.Fatalf("Post() error = %v", err)
	}
}
