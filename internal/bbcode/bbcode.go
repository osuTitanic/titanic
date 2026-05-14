package bbcode

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/Lekuruu/bbgo"
	"github.com/osuTitanic/titanic-go/internal/constants"
)

var (
	defaultRenderer = New(Options{})
	timecodeRegex   = regexp.MustCompile(`\b(?:osu://edit/)?(\d{2,}):([0-5]\d)[:.](\d{3})(?:\s(\((?:\d+[,|])*\d+\)))?`)
)

// Options contains bbcode rendering settings
type Options struct {
	BaseUrl string
}

// Renderer wraps bbgo with the tag set used by the application
type Renderer struct {
	parser *bbgo.BBGO
}

// New creates a renderer with titanic-specific tags registered
func New(options Options) *Renderer {
	parser := bbgo.NewWithOptions(bbgo.ParserOptions{
		UnknownLineRenderer: renderUnknownLine,
	})

	registerSimpleTags(parser)
	registerRawTags(parser)
	registerContainerTags(parser)
	registerLinkTags(parser, options)
	registerMediaTags(parser)

	return &Renderer{parser: parser}
}

// RenderHtml renders input using the default renderer
func RenderHtml(input string) string {
	return defaultRenderer.RenderHtml(input)
}

// Strip removes recognized bbcode tags using the default renderer
func Strip(input string, stripNewlines bool) string {
	return defaultRenderer.Strip(input, stripNewlines)
}

// RenderHtml renders input as html
func (r *Renderer) RenderHtml(input string) string {
	return renderTimecodes(r.parser.Parse(input))
}

// Strip removes recognized bbcode tags
func (r *Renderer) Strip(input string, stripNewlines bool) string {
	return r.parser.Strip(input, stripNewlines)
}

func registerSimpleTags(parser *bbgo.BBGO) {
	options := embeddedOptions()
	parser.AddSimpleFormatter("b", "<b>%s</b>", options)
	parser.AddSimpleFormatter("i", "<i>%s</i>", options)
	parser.AddSimpleFormatter("u", "<u>%s</u>", options)
	parser.AddSimpleFormatter("heading", "<h2>%s</h2>", options)
	parser.AddSimpleFormatter("strike", "<strike>%s</strike>", options)
	parser.AddSimpleFormatter("centre", "<center>%s</center>", options)
}

func registerRawTags(parser *bbgo.BBGO) {
	options := rawOptions()
	parser.AddSimpleFormatter("code", "%s", options)
	parser.AddSimpleFormatter("c", "%s", options)
}

func registerContainerTags(parser *bbgo.BBGO) {
	parser.AddSimpleFormatter("*", "<li>%s</li>", sameTagClosesOptions())
	parser.AddSimpleFormatter("spoilerbox", `<div class="spoiler"><div class="spoiler-head" onclick="return toggleSpoiler(this);">SPOILER</div><div class="spoiler-body">%s</div></div>`, embeddedOptions())

	parser.AddFormatter("box", renderBox, embeddedOptions())
	parser.AddFormatter("color", renderColor, embeddedOptions())
	parser.AddFormatter("quote", renderQuote, embeddedOptions())
	parser.AddFormatter("size", renderSize, embeddedOptions())
	parser.AddFormatter("list", renderList, embeddedOptions())
}

func registerLinkTags(parser *bbgo.BBGO, options Options) {
	parser.AddFormatter("url", renderURL, urlOptions())
	parser.AddFormatter("google", renderGoogle, rawOptions())
	parser.AddFormatter("email", renderEmail, rawOptions())
	parser.AddFormatter("profile", renderProfile(options), embeddedOptions())
}

func registerMediaTags(parser *bbgo.BBGO) {
	parser.AddFormatter("img", renderImage, rawOptions())
	parser.AddFormatter("video", renderVideo, rawOptions())
	parser.AddFormatter("youtube", renderYoutube, rawOptions())
}

func embeddedOptions() bbgo.TagOptions {
	return bbgo.TagOptions{
		RenderEmbedded:    true,
		TransformNewlines: true,
		EscapeHTML:        true,
		ReplaceLinks:      true,
		ReplaceCosmetic:   true,
	}
}

func rawOptions() bbgo.TagOptions {
	options := embeddedOptions()
	options.SameTagCloses = true
	options.RenderEmbedded = false
	options.ReplaceLinks = false
	options.ReplaceCosmetic = false
	return options
}

func sameTagClosesOptions() bbgo.TagOptions {
	options := embeddedOptions()
	options.SameTagCloses = true
	return options
}

func urlOptions() bbgo.TagOptions {
	options := embeddedOptions()
	options.ReplaceLinks = false
	return options
}

func renderUnknownLine(tagText string, context bbgo.Context) (string, bool) {
	return fmt.Sprintf(`<div class="beatmap-header">%s</div>`, sanitizeInput(tagText)), true
}

func renderBox(ctx bbgo.RenderContext) string {
	title := sanitizeInput(ctx.Options.Get("box"))
	return fmt.Sprintf(`<div class="spoiler"><div class="spoiler-head" onclick="return toggleSpoiler(this);">%s</div><div class="spoiler-body">%s</div></div>`, title, ctx.Value)
}

func renderColor(ctx bbgo.RenderContext) string {
	color := strings.ReplaceAll(sanitizeInput(ctx.Options.Get("color")), ";", "")
	return fmt.Sprintf(`<span style="color:%s;">%s</span>`, color, ctx.Value)
}

func renderQuote(ctx bbgo.RenderContext) string {
	author := ctx.Options.Get("quote")
	if author == "" {
		return fmt.Sprintf("<blockquote>%s</blockquote>", ctx.Value)
	}
	return fmt.Sprintf("<blockquote><h4>%s wrote:</h4><i>%s</i></blockquote>", sanitizeInput(author), ctx.Value)
}

func renderSize(ctx bbgo.RenderContext) string {
	size, ok := resolveSize(ctx.Options.Get("size"))
	if !ok {
		return ctx.Value
	}
	return fmt.Sprintf(`<span style="font-size:%d%%;">%s</span>`, size, ctx.Value)
}

func renderList(ctx bbgo.RenderContext) string {
	if _, ok := ctx.Options["list"]; ok {
		return fmt.Sprintf("<ol>%s</ol>", ctx.Value)
	}
	return fmt.Sprintf("<ul>%s</ul>", ctx.Value)
}

func renderURL(ctx bbgo.RenderContext) string {
	link := ctx.Options.Get("url")
	if link == "" {
		link = ctx.Value
	}
	link = sanitizeUrl(link)
	return fmt.Sprintf(`<a href="%s" target="_blank">%s</a>`, link, ctx.Value)
}

func renderGoogle(ctx bbgo.RenderContext) string {
	return fmt.Sprintf(`<a href="https://letmegooglethat.com/?q=%s" target="_blank">%s</a>`, ctx.Value, ctx.Value)
}

func renderEmail(ctx bbgo.RenderContext) string {
	email := ctx.Options.Get("email")
	if email == "" {
		email = ctx.Value
	}
	email = sanitizeInput(email)
	if !constants.Email.MatchString(email) {
		return ctx.Value
	}
	return fmt.Sprintf(`<a href="mailto:%s">%s</a>`, email, ctx.Value)
}

func renderProfile(options Options) bbgo.RenderFunc {
	return func(ctx bbgo.RenderContext) string {
		profile := ctx.Options.Get("profile")
		if profile == "" {
			profile = ctx.Value
		}
		return fmt.Sprintf(`<a href="%s/u/%s">%s</a>`, strings.TrimRight(options.BaseUrl, "/"), sanitizeInput(profile), ctx.Value)
	}
}

func renderImage(ctx bbgo.RenderContext) string {
	source := validAbsoluteUrl(ctx.Value)
	if source == "" {
		return ""
	}
	return fmt.Sprintf(`<img src="%s" loading="lazy">`, sanitizeInput(source))
}

func renderVideo(ctx bbgo.RenderContext) string {
	source := validAbsoluteUrl(ctx.Value)
	if source == "" {
		return ""
	}
	return fmt.Sprintf(`<video src="%s" controls></video>`, sanitizeInput(source))
}

func renderYoutube(ctx bbgo.RenderContext) string {
	// TODO: maybe make this string not 10 billion characters long lmao
	videoID := youtubeId(ctx.Value)
	return fmt.Sprintf(`<iframe width="373" height="210" src="https://www.youtube.com/embed/%s" title="YouTube Video Player" frameborder="0" allow="accelerometer; autoplay;clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>`, sanitizeInput(videoID))
}

func resolveSize(input string) (int, bool) {
	if size, ok := namedSizes[input]; ok {
		return size, true
	}

	size, err := strconv.Atoi(input)
	if err != nil {
		return 0, false
	}
	return max(10, min(800, size)), true
}

var namedSizes = map[string]int{
	"tiny":   50,
	"small":  85,
	"normal": 100,
	"large":  180,
}

func sanitizeInput(text string) string {
	replacements := []struct {
		find string
		with string
	}{
		{"&", "&amp;"},
		{"<", "&lt;"},
		{">", "&gt;"},
		{`"`, "&quot;"},
		{"'", "&#39;"},
	}

	for _, replacement := range replacements {
		text = strings.ReplaceAll(text, replacement.find, replacement.with)
	}
	return text
}

func sanitizeUrl(text string) string {
	text = sanitizeInput(text)
	if !strings.HasPrefix(text, "http") {
		text = "http://" + text
	}
	return text
}

func validAbsoluteUrl(input string) string {
	parsed, err := url.Parse(strings.TrimSpace(input))
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return ""
	}
	return input
}

func youtubeId(input string) string {
	input = strings.TrimSpace(input)
	if strings.Contains(input, "/") {
		parts := strings.Split(input, "/")
		input = parts[len(parts)-1]
	}
	return strings.ReplaceAll(input, "watch?v=", "")
}

func renderTimecodes(input string) string {
	return timecodeRegex.ReplaceAllStringFunc(input, func(match string) string {
		timecode := strings.TrimPrefix(match, "osu://edit/")
		return fmt.Sprintf(`<a class="beatmap-timecode" href="osu://edit/%s">%s</a>`, timecode, timecode)
	})
}
