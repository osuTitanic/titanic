package routes

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/repositories"
	"github.com/osuTitanic/titanic/services/stern/internal/server"
	"github.com/osuTitanic/titanic/services/stern/internal/templates"
)

func Search(ctx *server.Context) {
	query := ctx.Request.URL.Query()
	userId := new(int)
	if ctx.CurrentUser != nil {
		userId = &ctx.CurrentUser.Id
	}

	options, page := buildBeatmapsetSearchOptions(query, userId)
	result, err := ctx.State.Repositories.Beatmapsets.SearchPage(options, "Beatmaps")
	if err != nil {
		ctx.Logger.Error("Failed to search beatmapsets", "options", options, "error", err)
		InternalServerError(ctx)
		return
	}

	pagination := templates.NewPagination(templates.PaginationOptions{
		Path:        "/beatmapsets",
		Query:       query,
		CurrentPage: page,
		PageSize:    result.Options.Limit,
		Total:       int(result.Total),
	})
	view := templates.BeatmapSearchView{
		DefaultView: buildDefaultView(ctx),
		Beatmapsets: result.Beatmapsets,
		SearchSort:  strconv.Itoa(int(result.Options.Sort)),
		SearchOrder: strconv.Itoa(int(result.Options.Order)),
		Pagination:  pagination,
	}
	ctx.RenderTemplate(http.StatusOK, "pages/public/search", view)
}

func buildBeatmapsetSearchOptions(query url.Values, userId *int) (repositories.BeatmapsetSearchOptions, int) {
	options := repositories.BeatmapsetSearchOptions{
		QueryString: query.Get("query"),
		Order:       constants.SearchOrderDescending,
		Category:    constants.BeatmapCategoryLeaderboard,
		Sort:        constants.BeatmapSortRanked,
		Limit:       50,
	}

	if genre, ok := parseInt(query.Get("genre")); ok {
		value := constants.BeatmapGenre(genre)
		options.Genre = &value
	}
	if language, ok := parseInt(query.Get("language")); ok {
		value := constants.BeatmapLanguage(language)
		options.Language = &value
	}
	if category, ok := parseInt(query.Get("category")); ok {
		options.Category = constants.BeatmapCategory(category)
	}
	if sort, ok := parseInt(query.Get("sort")); ok {
		options.Sort = constants.BeatmapSort(sort)
	}
	if order, ok := parseInt(query.Get("order")); ok {
		options.Order = constants.SearchOrder(order)
	}
	if mode, ok := parseInt(query.Get("mode")); ok {
		value := constants.Mode(mode)
		options.Mode = &value
	}

	page := 1
	if parsed, ok := parseInt(query.Get("page")); ok && parsed > 1 {
		page = parsed
	}
	options.Offset = (page - 1) * options.Limit

	options.HasVideo = query.Get("video") != ""
	options.HasStoryboard = query.Get("storyboard") != ""
	options.TitanicOnly = query.Get("titanic") != ""

	if userId != nil {
		options.UserId = userId
		options.Played = query.Get("played") != ""
		options.Unplayed = query.Get("unplayed") != ""
		options.Cleared = query.Get("cleared") != ""
		options.Uncleared = query.Get("uncleared") != ""
	}
	return options, page
}

func parseInt(value string) (int, bool) {
	if value == "" {
		return 0, false
	}
	parsed, err := strconv.Atoi(value)
	return parsed, err == nil
}
