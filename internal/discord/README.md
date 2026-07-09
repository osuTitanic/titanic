# Discord

This module provides utility for posting messages to Discord webhooks. It supports plain JSON webhook payloads, embeds & multipart payloads with attached files.

## Usage

Create a webhook & call `Post()` on it. Note that webhooks must contain at least text content, a file, or an embed.

```go
content := "Hello from titanic-go"

webhook := discord.Webhook{
	URL:     "https://discord.com/api/webhooks/...",
	Content: &content,
}
if err := webhook.Post(); err != nil {
	return err
}
```

Use `discord.Embed` for uh... embeds.

```go
embed := discord.Embed{
	Title:       ptr("w00t p00t"),
	Description: ptr("ahoy!"),
	Color:       ptr(0xFF66AB),
}
embed.AddField("Server Status", "crashing", true)
embed.AddField("Server Uptime", "-4 seconds", true)

webhook := discord.Webhook{
	URL:      "https://discord.com/api/webhooks/...",
	Username: ptr("BanchoBot"),
}
webhook.AddEmbed(embed)

if err := webhook.Post(); err != nil {
	return err
}
```

Attach file data with `SetFile` or `SetFileReader` for streaming.

```go
webhook := discord.Webhook{
	URL: "https://discord.com/api/webhooks/...",
}
webhook.SetFile("message.txt", []byte("nyello!"))

// or ...
// webhook.SetFileReader("message.txt", bytes.NewReader([]byte("nyello!")))

if err := webhook.Post(); err != nil {
	return err
}
```

## Officer

Officer is a webhook dispatcher for staff notifications. Configure `OFFICER_WEBHOOK_URL` for the staff webhook.
Officer calls accept tags already, but all tags currently use this same webhook URL. In the future, we may support multiple webhooks for different tags, so that notifications go in different channels based on the tag.

Use the shared state instance from service code:

```go
err := app.Officer.Call(discord.OfficerTagScoreSubmission, "Suspicious score submitted")
if err != nil {
	return err
}
```

If `OFFICER_WEBHOOK_URL` is not configured, calls are ignored.

## Testing

The webhook tests are integration tests. They are skipped unless you provide a Discord webhook URL with the `-webhook-url` flag.

```sh
go test ./internal/discord -webhook-url "https://discord.com/api/webhooks/..."
```
