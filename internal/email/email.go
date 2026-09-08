package email

import (
	"errors"
	"fmt"
	"mime"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

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

// Validate ensures required fields exist & header injection is not a thing
func (message *Message) Validate() error {
	if message == nil {
		return errors.New("email: message is nil")
	}

	if len(message.To) == 0 {
		return errors.New("email: at least one recipient is required")
	}

	if message.Subject == "" {
		return errors.New("email: subject is required")
	}

	if message.TextBody == "" && message.HTMLBody == "" {
		return errors.New("email: text or HTML body is required")
	}

	for _, recipient := range message.To {
		if err := validateHeader("To", recipient); err != nil {
			return err
		}
	}
	for name, value := range message.Headers {
		if err := validateHeader(name, value); err != nil {
			return err
		}
	}
	if err := validateHeader("Subject", message.Subject); err != nil {
		return err
	}

	return nil
}

func (message *Message) BuildMimeMessage(from string) ([]byte, error) {
	if err := message.Validate(); err != nil {
		return nil, err
	}
	if err := validateHeader("From", from); err != nil {
		return nil, err
	}

	var builder strings.Builder
	encodedSubject := mime.QEncoding.Encode("utf-8", message.Subject)
	builder.WriteString("From: " + from + "\r\n")
	builder.WriteString("To: " + strings.Join(message.To, ", ") + "\r\n")
	builder.WriteString("Subject: " + encodedSubject + "\r\n")
	builder.WriteString("MIME-Version: 1.0\r\n")

	for header, value := range message.Headers {
		switch strings.ToLower(header) {
		case "from", "to", "subject", "date", "mime-version", "content-type", "content-transfer-encoding":
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

func validateHeader(name, value string) error {
	if name == "" {
		return errors.New("email: header name is required")
	}
	if !utf8.ValidString(value) {
		return fmt.Errorf("email: invalid UTF-8 in %s header", name)
	}
	for _, c := range name {
		if c < '!' || c > '~' || c == ':' {
			return errors.New("email: invalid header name")
		}
	}
	for _, c := range value {
		if c != '\t' && unicode.IsControl(c) {
			return fmt.Errorf("email: invalid control character in %s header", name)
		}
	}
	return nil
}
