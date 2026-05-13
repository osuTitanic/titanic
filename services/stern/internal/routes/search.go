package routes

import (
	"net/http"

	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
)

func Search(ctx *server.Context) {
	view := templates.BeatmapSearchView{
		DefaultView: buildDefaultView(ctx),
		Beatmapsets: nil, // TODO
	}
	ctx.RenderTemplate(http.StatusOK, "pages/public/search", view)
}
