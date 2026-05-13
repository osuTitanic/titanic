package routes

import (
	"net/http"

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

	view := templates.BeatmapSearchView{
		DefaultView: buildDefaultView(ctx),
		Beatmapsets: nil, // TODO
		SearchSort:  searchSort,
		SearchOrder: searchOrder,
	}
	ctx.RenderTemplate(http.StatusOK, "pages/public/search", view)
}
