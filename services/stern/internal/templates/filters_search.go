package templates

import (
	"reflect"
	"sort"

	"github.com/CloudyKit/jet/v6"
	"github.com/osuTitanic/titanic-go/internal/schemas"
)

func beatmapDifficultySort(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("orderedBeatmaps", 1, 1)

	beatmaps, ok := a.Get(0).Interface().([]*schemas.Beatmap)
	if !ok {
		return reflect.ValueOf([]*schemas.Beatmap{})
	}
	if len(beatmaps) == 0 {
		return reflect.ValueOf(beatmaps)
	}

	sort.SliceStable(beatmaps, func(i, j int) bool {
		if beatmaps[i].Mode == beatmaps[j].Mode {
			return beatmaps[i].Diff < beatmaps[j].Diff
		}
		return beatmaps[i].Mode < beatmaps[j].Mode
	})
	return reflect.ValueOf(beatmaps)
}

func beatmapRatingWidth(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("beatmapRatingWidth", 1, 1)

	ratingAverage, ok := a.Get(0).Interface().(float64)
	if !ok {
		return reflect.ValueOf(0)
	}

	width := 100 - (ratingAverage/10)*100
	width = max(0, min(width, 100))
	return reflect.ValueOf(width)
}
