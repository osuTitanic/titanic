package wiki

import (
	stdhtml "html"
	"strings"
	"unicode"
)

func replaceTocMarker(output, toc string) string {
	replacements := []string{
		"<p>[TOC]</p>\n",
		"<p>[TOC]</p>",
		"<p>[TOC]<br>\n</p>\n",
		"<p>[TOC]<br>\n</p>",
	}
	for _, marker := range replacements {
		output = strings.ReplaceAll(output, marker, toc)
	}
	return output
}

func collectTocHeadings(text string) []tocHeading {
	var headings []tocHeading
	lines := strings.SplitSeq(text, "\n")

	// For each line we want to check if it is a heading and if so,
	// we want to extract the level and text of the heading.
	for line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "#") {
			// Not a heading, skip...
			continue
		}

		level := 0

		// Count the number of '#' characters at the start of the line to determine the heading level
		for level < len(trimmed) && trimmed[level] == '#' {
			level++
		}

		// Heading level must be between 1 and 4, and there must be a space after the hashes
		if level < 1 || level > 4 || level >= len(trimmed) || trimmed[level] != ' ' {
			continue
		}

		text := strings.TrimSpace(trimmed[level:])
		text = strings.TrimSpace(strings.TrimRight(text, "#"))
		text = plainHeadingText(text)
		if text == "" {
			continue
		}

		headings = append(headings, tocHeading{
			Level: level,
			Text:  text,
			Slug:  headingSlug(text),
		})
	}

	if len(headings) > 0 && headings[0].Level == 1 {
		// If the first heading is level 1, we skip it to
		// avoid including the main title in the toc
		headings = headings[1:]
	}
	return headings
}

func renderToc(headings []tocHeading) string {
	if len(headings) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.WriteString(`<div class="toc"><span class="toctitle">Contents</span><ul>`)
	previousLevel := headings[0].Level

	for index, heading := range headings {
		levelDiff := heading.Level - previousLevel
		switch {
		case levelDiff > 0:
			builder.WriteString(strings.Repeat("<ul>", levelDiff))
		case levelDiff < 0:
			builder.WriteString(strings.Repeat("</li></ul>", -levelDiff))
			builder.WriteString("</li>")
		case index > 0:
			builder.WriteString("</li>")
		}

		builder.WriteString(`<li><a href="#`)
		builder.WriteString(stdhtml.EscapeString(heading.Slug))
		builder.WriteString(`">`)
		builder.WriteString(stdhtml.EscapeString(heading.Text))
		builder.WriteString(`</a>`)
		previousLevel = heading.Level
	}

	depth := previousLevel - headings[0].Level + 1
	builder.WriteString("</li>")
	builder.WriteString(strings.Repeat("</ul>", depth))
	builder.WriteString("</div>")
	return builder.String()
}

// e.g. `**[Link](url)**` -> `Link`
func plainHeadingText(text string) string {
	text = wikiLinkPattern.ReplaceAllStringFunc(text, func(value string) string {
		match := wikiLinkPattern.FindStringSubmatch(value)
		if len(match) > 2 && match[2] != "" {
			return match[2]
		}
		if len(match) > 1 {
			return match[1]
		}
		return value
	})
	text = markdownLink.ReplaceAllString(text, "$1")
	text = strings.Trim(text, "*_`")
	return strings.TrimSpace(text)
}

// e.g. "My Heading" -> "my-heading"
func headingSlug(text string) string {
	text = strings.ToLower(text)
	var builder strings.Builder
	lastDash := false

	for _, r := range text {
		switch {
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			builder.WriteRune(r)
			lastDash = false
		case unicode.IsSpace(r) || r == '_' || r == '-':
			if !lastDash {
				builder.WriteRune('-')
				lastDash = true
			}
		}
	}

	return strings.Trim(builder.String(), "-")
}
