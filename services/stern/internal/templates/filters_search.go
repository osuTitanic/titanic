package templates

import (
	"net/url"
	"reflect"
	"sort"
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/osuTitanic/titanic/internal/schemas"
)

type SearchHiddenInput struct {
	Name  string
	Value string
}

func searchParamUrl(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("searchParamUrl", 4, 4)

	current, _ := a.Get(0).Interface().(url.Values)
	path, _ := a.Get(1).Interface().(string)
	key, _ := a.Get(2).Interface().(string)
	value, _ := a.Get(3).Interface().(string)

	query := cloneQuery(current)
	query.Del("page")

	if value == "" {
		query.Del(key)
	} else {
		query.Set(key, value)
	}
	return reflect.ValueOf(searchUrl(path, query))
}

func searchFlagUrl(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("searchFlagUrl", 4, 4)

	current, _ := a.Get(0).Interface().(url.Values)
	path, _ := a.Get(1).Interface().(string)
	key, _ := a.Get(2).Interface().(string)
	removes, _ := a.Get(3).Interface().(string)

	query := cloneQuery(current)
	query.Del("page")
	if query.Get(key) != "" {
		query.Del(key)
	} else {
		query.Set(key, "1")
	}

	for remove := range strings.SplitSeq(removes, ",") {
		remove = strings.TrimSpace(remove)
		if remove != "" {
			query.Del(remove)
		}
	}
	return reflect.ValueOf(searchUrl(path, query))
}

func searchSortUrl(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("searchSortUrl", 3, 5)

	current, _ := a.Get(0).Interface().(url.Values)
	path, _ := a.Get(1).Interface().(string)
	sort, _ := a.Get(2).Interface().(string)

	defaultSort := "4"
	defaultOrder := "0"
	if a.NumOfArguments() > 3 {
		defaultSort, _ = a.Get(3).Interface().(string)
	}
	if a.NumOfArguments() > 4 {
		defaultOrder, _ = a.Get(4).Interface().(string)
	}

	return reflect.ValueOf(buildSearchSortUrl(current, path, sort, defaultSort, defaultOrder))
}

func buildSearchSortUrl(current url.Values, path, sort, defaultSort, defaultOrder string) string {
	currentSort := current.Get("sort")
	if currentSort == "" {
		currentSort = defaultSort
	}
	currentOrder := current.Get("order")
	if currentOrder == "" {
		currentOrder = defaultOrder
	}

	query := cloneQuery(current)
	query.Del("page")
	query.Set("sort", sort)

	if currentSort == sort {
		if currentOrder == "0" {
			query.Set("order", "1")
		} else {
			query.Set("order", "0")
		}
	} else {
		query.Set("order", defaultOrder)
	}
	return searchUrl(path, query)
}

func searchHiddenInputs(a jet.Arguments) reflect.Value {
	a.RequireNumOfArguments("searchHiddenInputs", 1, 2)

	inputs := make([]SearchHiddenInput, 0)
	excluded := map[string]bool{"query": true, "page": true}
	if a.NumOfArguments() > 1 {
		extra, _ := a.Get(1).Interface().(string)
		for name := range strings.SplitSeq(extra, ",") {
			name = strings.TrimSpace(name)
			if name != "" {
				excluded[name] = true
			}
		}
	}

	current, _ := a.Get(0).Interface().(url.Values)
	for name, values := range current {
		if excluded[name] {
			continue
		}
		for _, value := range values {
			inputs = append(inputs, SearchHiddenInput{Name: name, Value: value})
		}
	}

	// Sort the inputs by name and then by value for consistent ordering
	sort.Slice(inputs, func(i, j int) bool {
		if inputs[i].Name == inputs[j].Name {
			return inputs[i].Value < inputs[j].Value
		}
		return inputs[i].Name < inputs[j].Name
	})
	return reflect.ValueOf(inputs)
}

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

func searchUrl(path string, query url.Values) string {
	if path == "" {
		path = "/"
	}
	if encoded := query.Encode(); encoded != "" {
		return path + "?" + encoded
	}
	return path
}
