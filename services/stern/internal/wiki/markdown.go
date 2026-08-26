package wiki

import (
	"bytes"
	stdhtml "html"
	"net/url"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	gmhtml "github.com/yuin/goldmark/renderer/html"
	gmtext "github.com/yuin/goldmark/text"
)

var (
	markdownRenderer = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			extension.DefinitionList,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			gmhtml.WithHardWraps(),
			gmhtml.WithUnsafe(),
		),
	)
	fencedCodeBlock = regexp.MustCompile(`(?s)<pre><code([^>]*)>(.*?)</code></pre>`)
	markdownLink    = regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`)
)

type tocHeading struct {
	Level int
	Text  string
	Slug  string
}

// RenderMarkdown converts the given Markdown text to HTML.
// It handles wiki-style links, table of contents, and code blocks for syntax highlighting.
func RenderMarkdown(text, language string) (string, error) {
	return renderMarkdown(text, language, "")
}

// RenderMarkdownWithSourceURL converts the given Markdown text to HTML and
// resolves relative image paths against the URL of the Markdown source file.
func RenderMarkdownWithSourceURL(text, language, sourceURL string) (string, error) {
	return renderMarkdown(text, language, sourceURL)
}

func renderMarkdown(text, language, sourceURL string) (string, error) {
	text = createValidUtf8(SanitizeMarkdown(text))
	text = StripFrontMatter(text)

	headings := collectTocHeadings(text)
	text = renderWikiLinks(text, language)

	source := []byte(text)
	document := markdownRenderer.Parser().Parse(gmtext.NewReader(source))
	resolveRelativeImageUrls(document, sourceURL)

	var buffer bytes.Buffer
	if err := markdownRenderer.Renderer().Render(&buffer, source, document); err != nil {
		return "", err
	}

	output := wrapCodeBlocks(buffer.String())
	output = replaceTocMarker(output, renderToc(headings))
	return output, nil
}

func resolveRelativeImageUrls(document ast.Node, sourceURL string) {
	base, err := url.Parse(sourceURL)
	if err != nil || !base.IsAbs() {
		return
	}

	ast.Walk(document, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		image, ok := node.(*ast.Image)
		if !ok {
			return ast.WalkContinue, nil
		}

		destination := strings.TrimSpace(string(image.Destination))
		reference, err := url.Parse(destination)
		if err != nil || !isRelativeImageUrl(reference) {
			return ast.WalkContinue, nil
		}

		image.Destination = []byte(base.ResolveReference(reference).String())
		return ast.WalkContinue, nil
	})
}

func isRelativeImageUrl(reference *url.URL) bool {
	return reference != nil &&
		!reference.IsAbs() &&
		reference.Host == "" &&
		reference.Path != "" &&
		!strings.HasPrefix(reference.Path, "/")
}

func renderWikiLinks(text, language string) string {
	var builder strings.Builder
	inFence := false

	for line := range strings.SplitAfterSeq(text, "\n") {
		trimmed := strings.TrimSpace(line)

		// Ignore replacing wiki links for code fences
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			inFence = !inFence
			builder.WriteString(line)
			continue
		}
		if inFence {
			builder.WriteString(line)
			continue
		}

		// Replace wiki links for this line
		builder.WriteString(renderWikiLinksInLine(line, language))
	}

	return builder.String()
}

func renderWikiLinksInLine(line, language string) string {
	matches := wikiLinkPattern.FindAllStringSubmatchIndex(line, -1)
	if len(matches) == 0 {
		return line
	}

	// A match will consist of:
	// [0] - start index of the whole match
	// [1] - end index of the whole match
	// [2] - start index of the link text
	// [3] - end index of the link text
	// [4] - start index of the label (if present)
	// [5] - end index of the label (if present)

	var builder strings.Builder
	var lastEnd = 0

	for _, match := range matches {
		start, end := match[0], match[1]

		// Check if this is a markdown link, if so, skip it
		if end < len(line) && line[end] == '(' {
			builder.WriteString(line[lastEnd:end])
			lastEnd = end
			continue
		}

		// Extract the link & label from the match
		link := strings.TrimSpace(line[match[2]:match[3]])
		label := link
		if match[4] >= 0 && match[5] >= 0 {
			label = strings.TrimSpace(line[match[4]:match[5]])
		}

		// Create a link that redirects to the target page
		href := "/wiki/" + link
		if language = NormalizeLanguage(language); language != "" {
			href = "/wiki/" + language + "/" + link
		}

		builder.WriteString(line[lastEnd:start])
		builder.WriteString(`<a class="wikilink" href="`)
		builder.WriteString(stdhtml.EscapeString(href))
		builder.WriteString(`">`)
		builder.WriteString(stdhtml.EscapeString(label))
		builder.WriteString(`</a>`)
		lastEnd = end
	}

	builder.WriteString(line[lastEnd:])
	return builder.String()
}

func wrapCodeBlocks(output string) string {
	// Following: https://github.com/Python-Markdown/markdown/blob/master/markdown/extensions/codehilite.py
	return fencedCodeBlock.ReplaceAllString(output, `<div class="codehilite"><pre><code$1>$2</code></pre></div>`)
}
