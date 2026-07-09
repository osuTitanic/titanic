package routes

import (
	"net/http"

	"github.com/osuTitanic/titanic/services/stern/internal/server"
)

func AccountLogout(ctx *server.Context) {
	if err := ctx.Request.ParseForm(); err != nil {
		ctx.Logger.Warn("Failed to parse logout form", "error", err)
	}

	redirectTarget := sanitizeRedirectTarget(ctx.Request.FormValue("redirect"))
	if redirectTarget == "" {
		redirectTarget = "/"
	}

	if !ctx.IsAuthenticated() {
		ctx.ExpireSessionCookie()
		ctx.Redirect(http.StatusSeeOther, redirectTarget)
		return
	}

	ok, err := ctx.ValidateCSRF()
	if err != nil {
		ctx.Logger.Error("Failed to validate logout csrf token", "user_id", ctx.CurrentUser.Id, "error", err)
		InternalServerError(ctx)
		return
	}
	if !ok {
		ctx.Logger.Warn("Invalid CSRF token on logout attempt", "user_id", ctx.CurrentUser.Id)
		ctx.ExpireSessionCookie()
		ctx.Redirect(http.StatusSeeOther, redirectTarget)
		return
	}

	if err := ctx.DeleteCurrentSessionCookie(); err != nil {
		ctx.Logger.Warn("Failed to delete website session", "user_id", ctx.CurrentUser.Id, "error", err)
	}
	if err := ctx.DeleteCurrentCSRFToken(); err != nil {
		ctx.Logger.Warn("Failed to delete csrf token", "user_id", ctx.CurrentUser.Id, "error", err)
	}

	ctx.ExpireSessionCookie()
	ctx.Redirect(http.StatusSeeOther, redirectTarget)
}
