package routes

import (
	"net/http"

	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
)

func Events(ctx *server.Context) {
	ctx.RenderTemplate(
		http.StatusOK, "pages/public/events",
		buildDefaultView(ctx),
	)
}
