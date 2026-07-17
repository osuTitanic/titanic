package routes

import (
	"regexp"
	"slices"
	"strings"

	"github.com/osuTitanic/titanic/services/stern/internal/server"
)

var (
	forumSmileyUrlRegex = regexp.MustCompile(
		// Copied from bbgo, we need this to avoid replacing smileys inside URLs
		`(?im)\b((?:https?://|www\d{0,3}[.]|[a-z0-9.\-]+[.][a-z]{2,}/)(?:[^\s()<>]+|\([^\s()<>]+\))+(?:\([^\s()<>]+\)|[^\s` + "`" + `!()\[\]{};:'".,<>?]))`,
	)
	forumSmileyRawTags = map[string]bool{
		// We don't want to replace smileys inside these tags
		"c":       true,
		"code":    true,
		"email":   true,
		"google":  true,
		"img":     true,
		"smiley":  true,
		"url":     true,
		"video":   true,
		"youtube": true,
	}
	forumSmileyReplacer = strings.NewReplacer(
		":)", "[smiley]50[/smiley]",
		";)", "[smiley]59[/smiley]",
		":D", "[smiley]242[/smiley]",
		":o", "[smiley]57[/smiley]",
		">:(", "[smiley]51[/smiley]",
		"8-)", "[smiley]248[/smiley]",
		":(", "[smiley]56[/smiley]",
		":?", "[smiley]60[/smiley]",
		":x", "[smiley]51[/smiley]",
		":P", "[smiley]58[/smiley]",
		":|", "[smiley]52[/smiley]",
		":!:", "[smiley]260[/smiley]",
		":?:", "[smiley]266[/smiley]",
		":cry:", "[smiley]55[/smiley]",
		":lol:", "[smiley]262[/smiley]",
		":roll:", "[smiley]269[/smiley]",
		":idea:", "[smiley]261[/smiley]",
		":oops:", "[smiley]268[/smiley]",
		":shock:", "[smiley]75[/smiley]",
		":arrow:", "[smiley]247[/smiley]",
	)
)

func normalizeForumPostSmileys(input string) string {
	var output strings.Builder
	var rawTag string

	// We will rebuild the input string piece by piece,
	// while replacing smileys outside of raw tags (like [code] or [url])
	// e.g. :shock: will be replaced with [smiley]75[/smiley]

	for len(input) > 0 {
		tagStart := strings.IndexByte(input, '[')
		if tagStart == -1 {
			// No more tags -> we can just replace smileys in the rest of the text
			writeForumSmileyText(&output, input, rawTag != "")
			break
		}

		// We found a tag, we will replace smileys in the text before it
		writeForumSmileyText(&output, input[:tagStart], rawTag != "")
		input = input[tagStart:]

		tagEnd := strings.IndexByte(input, ']')
		if tagEnd == -1 {
			// No closing bracket -> we can just replace smileys in the rest of the text
			writeForumSmileyText(&output, input, rawTag != "")
			break
		}

		// We found a tag, we will add it to the output and check if it's a raw tag
		// If not, we will continue replacing smileys in the text after it

		tag := input[:tagEnd+1]
		name, closing := forumSmileyTagName(tag)
		output.WriteString(tag)
		input = input[tagEnd+1:]

		if rawTag == "" && !closing && forumSmileyRawTags[name] {
			rawTag = name
		} else if rawTag != "" && name == rawTag && closing {
			rawTag = ""
		}
	}

	return output.String()
}

func forumSmileysEnabled(ctx *server.Context) bool {
	values, submitted := ctx.Request.PostForm["enable-smilies"]
	if !submitted {
		// on by default
		return true
	}
	return slices.Contains(values, "1")
}

func writeForumSmileyText(output *strings.Builder, text string, raw bool) {
	if raw {
		output.WriteString(text)
		return
	}

	// We don't want to replace smileys inside URLs, so we find all URLs and
	// replace smileys in the text outside of them

	last := 0
	for _, match := range forumSmileyUrlRegex.FindAllStringIndex(text, -1) {
		output.WriteString(forumSmileyReplacer.Replace(text[last:match[0]]))
		output.WriteString(text[match[0]:match[1]])
		last = match[1]
	}
	output.WriteString(forumSmileyReplacer.Replace(text[last:]))
}

func forumSmileyTagName(tag string) (string, bool) {
	body := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(tag, "["), "]"))
	closing := strings.HasPrefix(body, "/")
	body = strings.TrimSpace(strings.TrimPrefix(body, "/"))
	if end := strings.IndexAny(body, "= \t\r\n"); end >= 0 {
		body = body[:end]
	}
	return strings.ToLower(body), closing
}
