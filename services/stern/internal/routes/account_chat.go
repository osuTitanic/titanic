package routes

import (
	"net/http"

	"github.com/osuTitanic/titanic/services/stern/internal/server"
)

func AccountChat(ctx *server.Context) {
	if !ctx.RequireLogin() {
		return
	}

	ctx.RenderTemplate(
		http.StatusOK, "pages/account/settings_chat",
		buildDefaultView(ctx),
	)
}
