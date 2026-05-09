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

type BeatmapWithCountItem struct {
	Beatmap *schemas.Beatmap
	Count   int
}

func homeRenderNewsText(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("homeRenderNewsText", 1, 1)
	post := a.Get(0).Interface().(schemas.ForumPost)

	content := post.RenderForNews(homeNewsIgnoredTags...)
	return reflect.ValueOf(content)
}

func homeIterateMostPlayed(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("homeIterateMostPlayed", 1, 1)
	mostPlayed := a.Get(0).Interface().(map[int]*schemas.Beatmap)

	items := make([]BeatmapWithCountItem, 0, len(mostPlayed))

	for count, beatmap := range mostPlayed {
		items = append(items, BeatmapWithCountItem{
			Beatmap: beatmap,
			Count:   count,
		})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Count > items[j].Count
	})

	return reflect.ValueOf(items)
}
