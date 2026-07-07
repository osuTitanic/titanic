package wiki

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	// wikiLinkPattern matches wiki-style links in the format [[path|label]] or [[path]].
	wikiLinkPattern = regexp.MustCompile(`\[\[([^|\]]+)(?:\|([^\]]+))?\]\]`)

	// frontMatter matches yaml front matter at the beginning of a markdown file,
	// which is typically enclosed between '---' lines.
	frontMatter = regexp.MustCompile(`(?s)^---\r?\n.*?\r?\n---\r?\n?`)
)

type Outlink struct {
	Path  string
	Label string
}

// PageName extracts the page name from a given path.
// e.g. "folder/subfolder/page.md" -> "page"
func PageName(path string) string {
	path = strings.TrimSuffix(strings.TrimSpace(path), "/")
	path = strings.TrimSuffix(path, ".md")

	if index := strings.LastIndex(path, "/"); index >= 0 {
		path = path[index+1:]
	}
	return strings.ReplaceAll(path, "_", " ")
}

// PagePath extracts the page path from a given path.
// e.g. "folder/subfolder/page.md" -> "folder/subfolder"
func PagePath(path string) string {
	path = strings.TrimSuffix(strings.TrimSpace(path), "/")
	path = strings.TrimSuffix(path, ".md")
	return strings.ReplaceAll(path, "_", " ")
}

// FormatPath formats a given path and page name into a standardized path format.
// e.g. FormatPath("folder/subfolder", "page") -> "folder/subfolder/page"
func FormatPath(path, pageName string) string {
	tree := strings.Split(strings.TrimSpace(strings.TrimSuffix(strings.ReplaceAll(path, "_", " "), "/")), "/")
	if len(tree) == 0 {
		return pageName
	}

	for index, title := range tree {
		tree[index] = strings.ReplaceAll(title, " ", "_")
	}
	tree[len(tree)-1] = pageName
	return strings.Join(tree, "/")
}

// SanitizeMarkdown cleans up the input markdown text.
func SanitizeMarkdown(text string) string {
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "\ufeff")
	text = strings.TrimPrefix(text, string([]byte{0xef, 0xbb, 0xbf}))
	return strings.Trim(text, "\n")
}

// StripFrontMatter removes yaml front matter from the beginning of a markdown text.
// e.g. "---\ntitle: Example\n---\nContent" -> "Content"
func StripFrontMatter(text string) string {
	return frontMatter.ReplaceAllString(text, "")
}

// ParseTitle extracts the title from the first line of the markdown text.
// e.g. "# Example Title\nContent" -> "Example Title"
func ParseTitle(text string) string {
	text = StripFrontMatter(text)
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return ""
	}
	return strings.TrimSpace(strings.TrimLeft(lines[0], "#"))
}

// ExtractOutlinks finds all wiki-style links in the given
// markdown text & returns them as a slice of `Outlink`.
func ExtractOutlinks(text string) []Outlink {
	matches := wikiLinkPattern.FindAllStringSubmatch(text, -1)
	outlinks := make([]Outlink, 0, len(matches))

	// match[1] is the path, match[2] is the optional label
	// so e.g. [[path|label]] -> match[1] = "path", match[2] = "label"
	for _, match := range matches {
		path := strings.TrimSpace(match[1])
		if path == "" {
			continue
		}

		label := path
		if len(match) > 2 && match[2] != "" {
			label = strings.TrimSpace(match[2])
		}
		outlinks = append(outlinks, Outlink{Path: path, Label: label})
	}
	return outlinks
}

func createValidUtf8(text string) string {
	if utf8.ValidString(text) {
		return text
	}
	return strings.ToValidUTF8(text, "")
}
