# BBCode

This module renders BBCode used by Titanic services. It utilizes [bbgo](https://github.com/Lekuruu/bbgo) & adds new tags, media proxy stuff, and beatmap timecode rendering.

Forum posts use the renderer through `bbcode.RenderHtml(...)`. The website template engine configures that default renderer during startup.

## Usage

Render BBCode into HTML:

```go
html := bbcode.RenderHtml(post.Content)
```

Strip recognized BBCode tags when plain text is required:

```go
text := bbcode.Strip(post.Content, false)
```

Configure the default renderer once during startup with options from the config:

```go
bbcode.ConfigureDefault(bbcode.Options{
	BaseUrl:            cfg.OsuBaseUrl(),
	ValidImageServices: cfg.ValidImageServices(),
	ImageProxyBaseUrl:  cfg.ImageProxyBaseUrl,
	ImageProxySecret:   cfg.FrontendSecretKey,
})
```

Create a non-default renderer when tests or tooling need their own options:

```go
renderer := bbcode.New(bbcode.Options{
	BaseUrl: "https://osu.titanic.sh",
})

html := renderer.RenderHtml("[profile=2]Levi[/profile]")
```

## Supported Tags

The renderer currently supports:

- basic formatting: `b`, `i`, `u`, `strike`, `heading`, `centre`
- containers: `quote`, `box`, `spoilerbox`, `color`, `size`, `list`, `*`
- raw text: `code`, `c`
- links: `url`, `email`, `google`, `profile`
- media: `img`, `smiley`, `video`, `youtube`

Unknown line tags, e.g. `[Header]`, are rendered as beatmap headers. Beatmap timecodes such as `01:23:456` and `osu://edit/01:23:456` are converted into `osu://edit/...` links.

## Media Proxy

`img` and `video` tags support willnorris's [imageproxy](https://github.com/willnorris/imageproxy) service. When `ImageProxyBaseUrl` is set, media from hosts outside `ValidImageServices` is modified to a signed URL using `ImageProxySecret`.
