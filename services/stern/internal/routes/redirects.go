package routes

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/services/stern/internal/server"
)

func LegacyPageRedirect(ctx *server.Context) {
	page := ctx.PathValue("page")
	query := ctx.Request.URL.Query()

	switch page {
	case "download":
		ctx.Redirect(http.StatusFound, "/download")
	case "team":
		ctx.Redirect(http.StatusFound, "/g/1")
	case "pp", "ranking":
		ctx.Redirect(http.StatusFound, "/rankings/osu/performance")
	case "countryranking":
		ctx.Redirect(http.StatusFound, "/rankings/osu/country")
	case "player":
		if username := query.Get("f"); username != "" {
			ctx.Redirect(http.StatusFound, "/u/"+url.PathEscape(username))
			return
		}
		ctx.Redirect(http.StatusFound, "/")
	case "profile":
		if id := query.Get("u"); id != "" {
			ctx.Redirect(http.StatusFound, "/u/"+url.PathEscape(id))
			return
		}
		ctx.Redirect(http.StatusFound, "/")
	case "beatmap", "song":
		if id := query.Get("b"); id != "" {
			ctx.Redirect(http.StatusFound, "/b/"+url.PathEscape(id))
			return
		}
		if id := query.Get("s"); id != "" {
			ctx.Redirect(http.StatusFound, "/s/"+url.PathEscape(id))
			return
		}
		ctx.Redirect(http.StatusFound, "/")
	case "playerranking":
		ctx.Redirect(http.StatusFound, playerRankingLocation(ctx))
	case "beatmaplist":
		ctx.Redirect(http.StatusFound, beatmapListLocation(query))
	default:
		ctx.Redirect(http.StatusFound, "/")
	}
}

func playerRankingLocation(ctx *server.Context) string {
	query := ctx.Request.URL.Query()
	mode := constants.ModeOsu
	if value, ok := parseInt(query.Get("m")); ok && value >= 0 && value <= 3 {
		mode = constants.Mode(value)
	}

	arguments := make([]string, 0, 3)
	if value, ok := queryValue(query, "f"); ok {
		arguments = append(arguments, "jumpto="+url.QueryEscape(value))
	}
	if value, ok := queryValue(query, "c"); ok {
		arguments = append(arguments, "country="+url.QueryEscape(value))
	}
	if value, ok := parseInt(query.Get("page")); ok {
		arguments = append(arguments, "page="+strconv.Itoa(value))
	}

	location := "/rankings/" + mode.Alias() + "/performance"
	if len(arguments) > 0 {
		location += "?" + strings.Join(arguments, "&")
	}
	if fragment := ctx.Request.URL.Fragment; fragment != "" {
		location += "#" + fragment
	}
	return location
}

func beatmapListLocation(query url.Values) string {
	arguments := make([]string, 0, 6)

	for _, argument := range []struct {
		legacy string
		modern string
	}{
		{legacy: "la", modern: "language"},
		{legacy: "g", modern: "genre"},
		{legacy: "m", modern: "mode"},
		{legacy: "s", modern: "sort"},
		{legacy: "o", modern: "order"},
		{legacy: "q", modern: "query"},
	} {
		value, ok := queryValue(query, argument.legacy)
		if !ok {
			continue
		}

		if argument.legacy != "q" {
			if parsed, ok := parseInt(value); ok {
				value = strconv.Itoa(parsed)
			} else {
				continue
			}
		}
		arguments = append(arguments, argument.modern+"="+url.QueryEscape(value))
	}

	if len(arguments) == 0 {
		return "/beatmapsets"
	}
	return "/beatmapsets?" + strings.Join(arguments, "&")
}

func queryValue(query url.Values, key string) (string, bool) {
	values, ok := query[key]
	if !ok || len(values) == 0 {
		return "", false
	}
	return values[0], true
}
