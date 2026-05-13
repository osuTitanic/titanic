package routes

import (
	"net/http"

	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
)

func Search(ctx *server.Context) {
	query := ctx.Request.URL.Query()
	searchSort := query.Get("sort")
	if searchSort == "" {
		searchSort = "4" // Ranked
	}

	searchOrder := query.Get("order")
	if searchOrder == "" {
		searchOrder = "0" // Descending
	}
	// TODO: Move order & sort to enums

	testBeatmapset, err := ctx.State.Repositories.Beatmapsets.ById(292301, "Beatmaps")
	if err != nil {
		InternalServerError(ctx)
		return
	}
	// TODO: Implement actual search functionality

	view := templates.BeatmapSearchView{
		DefaultView: buildDefaultView(ctx),
		Beatmapsets: []*schemas.Beatmapset{testBeatmapset},
		SearchSort:  searchSort,
		SearchOrder: searchOrder,
	}
	ctx.RenderTemplate(http.StatusOK, "pages/public/search", view)
}
