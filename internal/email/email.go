package email

import (
	"errors"
	"fmt"
	"mime"
	"strings"
	"time"

	"github.com/osuTitanic/titanic/internal/config"
)

// Email defines the contract for email delivery backends
type Email interface {
	Setup() error
	Send(message *Message) error
	FromAddress() string
}

// NewEmailFromConfig constructs an Email implementation based on the provided type
func NewEmailFromConfig(config *config.Config) (Email, error) {
	switch config.EmailProvider {
	case "", "noop":
		return NewNoopEmail(config.EmailSender), nil
	case "smtp":
		return NewSMTPEmail(config), nil
	default:
		return nil, fmt.Errorf("email: unsupported email type: %s", config.EmailProvider)
	}
}

// Message describes an outbound email payload
type Message struct {
	To       []string
	Subject  string
	TextBody string
	HTMLBody string
	Headers  map[string]string
}

// Validate ensures the message has the required fields populated
func (m *Message) Validate() error {
	if m == nil {
		return errors.New("email: message is nil")
	}

	if len(m.To) == 0 {
		return errors.New("email: at least one recipient is required")
	}

	if m.Subject == "" {
		return errors.New("email: subject is required")
	}

	if m.TextBody == "" && m.HTMLBody == "" {
		return errors.New("email: text or HTML body is required")
	}

	return nil
}

func (message *Message) BuildMimeMessage(from string) ([]byte, error) {
	if len(message.To) == 0 {
		return nil, errors.New("email: no recipients provided")
	}

	var builder strings.Builder
	encodedSubject := mime.QEncoding.Encode("utf-8", message.Subject)
	builder.WriteString("From: " + from + "\r\n")
	builder.WriteString("To: " + strings.Join(message.To, ", ") + "\r\n")
	builder.WriteString("Subject: " + encodedSubject + "\r\n")
	builder.WriteString("MIME-Version: 1.0\r\n")

	for header, value := range message.Headers {
		if header == "From" || header == "To" || header == "Subject" {
			continue
		}
		builder.WriteString(header + ": " + value + "\r\n")
	}

	builder.WriteString("Date: " + time.Now().UTC().Format(time.RFC1123Z) + "\r\n")

	switch {
	case message.TextBody != "" && message.HTMLBody != "":
		// Create a multipart/alternative message
		boundary := fmt.Sprintf("puush-%d", time.Now().UnixNano())
		builder.WriteString("Content-Type: multipart/alternative; boundary=\"" + boundary + "\"\r\n\r\n")
		builder.WriteString("--" + boundary + "\r\n")
		builder.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		builder.WriteString("Content-Transfer-Encoding: 8bit\r\n\r\n")
		builder.WriteString(message.TextBody + "\r\n")
		builder.WriteString("--" + boundary + "\r\n")
		builder.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		builder.WriteString("Content-Transfer-Encoding: 8bit\r\n\r\n")
		builder.WriteString(message.HTMLBody + "\r\n")
		builder.WriteString("--" + boundary + "--\r\n")
	case message.HTMLBody != "":
		// HTML only message
		builder.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		builder.WriteString("Content-Transfer-Encoding: 8bit\r\n\r\n")
		builder.WriteString(message.HTMLBody)
	default:
		// Text only message
		builder.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		builder.WriteString("Content-Transfer-Encoding: 8bit\r\n\r\n")
		builder.WriteString(message.TextBody)
	}

	builder.WriteString("\r\n")
	return []byte(builder.String()), nil
}
