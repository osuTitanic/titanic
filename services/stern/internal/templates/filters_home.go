package templates

import (
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/osuTitanic/titanic-go/internal/schemas"
)

type HomeMostPlayedRow struct {
	PlayCount int
	Beatmap   schemas.Beatmap
}

var homeNewsIgnoredTags = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\[/?b\]`),
	regexp.MustCompile(`(?i)\[/?centre\]`),
	regexp.MustCompile(`(?i)\[size(?:=[^\]]+)?\]`),
	regexp.MustCompile(`(?i)\[/size\]`),
}

func homeRenderNewsText(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("homeRenderNewsText", 1, 1)
	post := a.Get(0).Interface().(schemas.ForumPost)

	for line := range strings.SplitSeq(post.Content, "\n") {
		lowercaseLine := strings.ToLower(line)
		if strings.Contains(lowercaseLine, "[heading]") || strings.Contains(lowercaseLine, "[img]") {
			continue
		}

		content := strings.TrimSpace(line)
		for _, regex := range homeNewsIgnoredTags {
			content = regex.ReplaceAllString(content, "")
		}

		// TODO: BBCode rendering
		content = strings.TrimSpace(content)
		if content != "" {
			return reflect.ValueOf(content)
		}
	}
	return reflect.ValueOf("")
}

func homeMostPlayedRows(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("homeMostPlayedRows", 1, 1)

	beatmaps := a.Get(0).Interface().(map[int]schemas.Beatmap)
	rows := make([]HomeMostPlayedRow, 0, len(beatmaps))

	for playCount, beatmap := range beatmaps {
		rows = append(rows, HomeMostPlayedRow{
			PlayCount: playCount,
			Beatmap:   beatmap,
		})
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].PlayCount > rows[j].PlayCount
	})
	return reflect.ValueOf(rows)
}
