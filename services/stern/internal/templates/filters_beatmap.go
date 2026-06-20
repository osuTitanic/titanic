package templates

import (
	"reflect"

	"github.com/CloudyKit/jet/v6"
	"github.com/osuTitanic/titanic-go/internal/constants"
)

func beatmapGenres(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("beatmapGenres", 0, 0)
	return reflect.ValueOf(constants.BeatmapGenres)
}

func beatmapLanguages(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("beatmapLanguages", 0, 0)
	return reflect.ValueOf(constants.BeatmapLanguages)
}

func beatmapStatuses(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("beatmapStatuses", 0, 0)
	return reflect.ValueOf(constants.BeatmapStatuses)
}
