package email

import (
	"io"
	"net/mail"
	"strings"
	"testing"
)

func TestMessageValidate(t *testing.T) {
	tests := []struct {
		name string
		msg  *Message
		want string
	}{
		{
			name: "nil message",
			want: "email: message is nil",
		},
		{
			name: "missing body",
			msg:  &Message{To: []string{"user@example.com"}, Subject: "Hello"},
			want: "email: text or HTML body is required",
		},
		{
			name: "valid",
			msg:  &Message{To: []string{"user@example.com"}, Subject: "Hello", TextBody: "Body"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.msg.Validate()
			if test.want == "" {
				if err != nil {
					t.Fatalf("Validate() error = %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("Validate() error = nil, want %q", test.want)
			}
			if err.Error() != test.want {
				t.Fatalf("Validate() error = %q, want %q", err.Error(), test.want)
			}
		})
	}
}

func TestBuildMimeMessage(t *testing.T) {
	msg := &Message{
		To:       []string{"user@example.com"},
		Subject:  "Hello",
		TextBody: "Plain body",
	}

	payload, err := msg.BuildMimeMessage("support@example.com")
	if err != nil {
		t.Fatalf("BuildMimeMessage() error = %v", err)
	}

	parsed, err := mail.ReadMessage(strings.NewReader(string(payload)))
	if err != nil {
		t.Fatalf("failed to parse message: %v", err)
	}

	if got := parsed.Header.Get("From"); got != "support@example.com" {
		t.Fatalf("From = %q, want %q", got, "support@example.com")
	}
	if got := parsed.Header.Get("To"); got != "user@example.com" {
		t.Fatalf("To = %q, want %q", got, "user@example.com")
	}
	if got := parsed.Header.Get("Subject"); got != "Hello" {
		t.Fatalf("Subject = %q, want %q", got, "Hello")
	}

	body, err := io.ReadAll(parsed.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	if string(body) != "Plain body\r\n" {
		t.Fatalf("body = %q, want %q", string(body), "Plain body\r\n")
	}
}
