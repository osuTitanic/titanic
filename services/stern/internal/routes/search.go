package routes

import (
	"net/http"
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

	options, page := buildBeatmapsetSearchOptions(ctx, userId)
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

func buildBeatmapsetSearchOptions(ctx *server.Context, userId *int) (repositories.BeatmapsetSearchOptions, int) {
	options := repositories.BeatmapsetSearchOptions{
		QueryString: ctx.QueryValue("query"),
		Order:       constants.SearchOrderDescending,
		Category:    constants.BeatmapCategoryLeaderboard,
		Sort:        constants.BeatmapSortRanked,
		Limit:       50,
	}

	if genre, err := ctx.QueryValueEnum[constants.BeatmapGenre]("genre"); err == nil {
		options.Genre = new(genre)
	}
	if language, err := ctx.QueryValueEnum[constants.BeatmapLanguage]("language"); err == nil {
		options.Language = new(language)
	}
	if category, err := ctx.QueryValueEnum[constants.BeatmapCategory]("category"); err == nil {
		options.Category = category
	}
	if sort, err := ctx.QueryValueEnum[constants.BeatmapSort]("sort"); err == nil {
		options.Sort = sort
	}
	if order, err := ctx.QueryValueEnum[constants.SearchOrder]("order"); err == nil {
		options.Order = order
	}
	if mode, err := ctx.QueryValueEnum[constants.Mode]("mode"); err == nil {
		options.Mode = new(mode)
	}

	page := 1
	if parsed, err := ctx.QueryValueInt("page"); err == nil && parsed > 1 {
		page = parsed
	}
	options.Offset = (page - 1) * options.Limit

	options.HasVideo = ctx.QueryValue("video") != ""
	options.HasStoryboard = ctx.QueryValue("storyboard") != ""
	options.TitanicOnly = ctx.QueryValue("titanic") != ""

	if userId != nil {
		options.UserId = userId
		options.Played = ctx.QueryValue("played") != ""
		options.Unplayed = ctx.QueryValue("unplayed") != ""
		options.Cleared = ctx.QueryValue("cleared") != ""
		options.Uncleared = ctx.QueryValue("uncleared") != ""
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
