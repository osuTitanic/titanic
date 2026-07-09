# Email

This module provides email functionality. It supports dummy delivery for development & smtp delivery for environments that should send real messages. It is built in such a way that additional providers can be added in the future.

Services should use the *state system* to access the email sender instance (`state.Email`). A manual setup is usually not required.

## Usage

Create an email sender from the app config, call `Setup()`, then send a `Message`.

```go
cfg, err := config.LoadConfig()
if err != nil {
	panic(err)
}

sender, err := email.NewEmailFromConfig(cfg)
if err != nil {
	return err
}
if err := sender.Setup(); err != nil {
	return err
}
// You may also use `state.Email` if you have the app state set up

message := &email.Message{
	To:       []string{"user@example.com"},
	Subject:  "yo",
	TextBody: "hi",
}
if err := sender.Send(message); err != nil {
	return err
}
```

Use `HTMLBody` for HTML email, or provide both `TextBody` and `HTMLBody` to build a multipart alternative message.

```go
message := &email.Message{
	To:       []string{"user@example.com"},
	Subject:  "something to put here",
	TextBody: "Open this email in an HTML client.",
	HTMLBody: "<p>HTML content</p>",
	Headers: map[string]string{
		"Reply-To": "support@titanic.sh",
	},
}
```

## Configuration

Email configuration is controlled by the email and smtp fields in `config.Config`.

- `EMAIL_PROVIDER` selects the delivery backend, defaulting to `noop`
- `EMAIL_SENDER` sets the default sender address
- `SMTP_HOST` sets the SMTP server hostname
- `SMTP_PORT` sets the SMTP server port, defaulting to `587`
- `SMTP_USER` and `SMTP_PASSWORD` enable SMTP auth when a user is provided
- `SMTP_TLS` enables STARTTLS when the server supports it
- `SMTP_SKIP_TLS_VERIFY` skips SMTP TLS certificate verification
