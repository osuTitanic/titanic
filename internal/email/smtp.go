package email

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/smtp"

	"github.com/osuTitanic/titanic/internal/config"
)

// SMTPEmail delivers messages using an SMTP server
type SMTPEmail struct {
	config *config.Config
	logger *slog.Logger
	auth   smtp.Auth
}

// NewSMTPEmail constructs an SMTP-backed email sender
func NewSMTPEmail(config *config.Config) Email {
	return &SMTPEmail{
		config: config,
		logger: slog.Default().With("component", "email"),
	}
}

// FromAddress returns the configured default sender address.
func (s *SMTPEmail) FromAddress() string {
	return s.config.EmailSender
}

// Setup validates the SMTP configuration and prepares any required auth
func (s *SMTPEmail) Setup() error {
	if s.config.SmtpHost == "" {
		return errors.New("email: SMTP host is required")
	}

	if s.config.SmtpPort == 0 {
		s.config.SmtpPort = 587
	}

	if s.config.EmailSender == "" {
		return errors.New("email: SMTP from address is required")
	}

	if s.config.SmtpUser != "" {
		s.auth = smtp.PlainAuth("", s.config.SmtpUser, s.config.SmtpPassword, s.config.SmtpHost)
	}

	return nil
}

// Send delivers the provided message using SMTP
func (s *SMTPEmail) Send(message *Message) error {
	if err := message.Validate(); err != nil {
		return err
	}

	mimeMessage, err := message.BuildMimeMessage(s.config.EmailSender)
	if err != nil {
		return err
	}

	address := fmt.Sprintf("%s:%d", s.config.SmtpHost, s.config.SmtpPort)
	client, err := smtp.Dial(address)
	if err != nil {
		return fmt.Errorf("email: failed to connect to SMTP server: %w", err)
	}
	defer client.Quit()

	if s.config.SmtpTls {
		if ok, _ := client.Extension("STARTTLS"); ok {
			tlsConfig := &tls.Config{
				ServerName:         s.config.SmtpHost,
				InsecureSkipVerify: s.config.SmtpSkipTlsVerify,
			}
			if err := client.StartTLS(tlsConfig); err != nil {
				return fmt.Errorf("email: failed to start TLS: %w", err)
			}
		}
	}

	if s.auth != nil {
		if ok, _ := client.Extension("AUTH"); ok {
			if err := client.Auth(s.auth); err != nil {
				return fmt.Errorf("email: failed to authenticate: %w", err)
			}
		}
	}

	if err := client.Mail(s.config.EmailSender); err != nil {
		return fmt.Errorf("email: failed to set sender: %w", err)
	}

	for _, recipient := range message.To {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("email: failed to add recipient %s: %w", recipient, err)
		}
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("email: failed to open SMTP data writer: %w", err)
	}
	defer writer.Close()

	if _, err := writer.Write(mimeMessage); err != nil {
		return fmt.Errorf("email: failed to write message body: %w", err)
	}

	s.logger.Debug("SMTP email sent", "to", message.To, "subject", message.Subject)
	return nil
}
