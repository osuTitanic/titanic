package templates

import (
	"reflect"
	"regexp"
	"sort"

	"github.com/CloudyKit/jet/v6"
	"github.com/osuTitanic/titanic-go/internal/schemas"
)

var homeNewsIgnoredTags = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\[/?b\]`),
	regexp.MustCompile(`(?i)\[/?centre\]`),
	regexp.MustCompile(`(?i)\[size(?:=[^\]]+)?\]`),
	regexp.MustCompile(`(?i)\[/size\]`),
	regexp.MustCompile(`(?i)\[/?heading\]`),
	regexp.MustCompile(`(?i)\[/?img\]`),
}

func homeRenderNewsText(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("homeRenderNewsText", 1, 1)
	post := a.Get(0).Interface().(schemas.ForumPost)

	content := post.RenderForNews(homeNewsIgnoredTags...)
	return reflect.ValueOf(content)
}

type HomeMostPlayedRow struct {
	PlayCount int
	Beatmap   schemas.Beatmap
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
