package routes

import (
	"net/http"

	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
)

func RenderResetPage(ctx *server.Context, errorMessage string) {
	view := templates.ResetView{
		DefaultView:  buildDefaultView(ctx),
		ErrorMessage: errorMessage,
	}
	ctx.RenderTemplate(http.StatusOK, "pages/account/reset", view)
}

func AccountPasswordReset(ctx *server.Context) {
	if ctx.IsAuthenticated() {
		ctx.Redirect(http.StatusSeeOther, "/")
		return
	}
	RenderResetPage(ctx, "")
}
