package wiki

import (
	"bytes"
	stdhtml "html"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	gmhtml "github.com/yuin/goldmark/renderer/html"
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
	text = createValidUtf8(SanitizeMarkdown(text))
	text = StripFrontMatter(text)

	headings := collectTocHeadings(text)
	text = renderWikiLinks(text, language)

	var buffer bytes.Buffer
	if err := markdownRenderer.Convert([]byte(text), &buffer); err != nil {
		return "", err
	}

	output := wrapCodeBlocks(buffer.String())
	output = replaceTocMarker(output, renderToc(headings))
	return output, nil
}

func renderWikiLinks(text, language string) string {
	var builder strings.Builder
	lines := strings.SplitAfter(text, "\n")
	inFence := false

	for _, line := range lines {
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
