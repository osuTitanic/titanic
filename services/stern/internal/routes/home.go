package routes

import (
	"net/http"

	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
)

func Home(ctx *server.Context) {
	view := templates.HomeView{
		DefaultView: BuildDefaultView(ctx),
		// TODO: Add more data
	}
	ctx.RenderTemplate(http.StatusOK, "pages/home", view)
}
