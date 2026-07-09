package routes

import (
	"net/http"

	"github.com/osuTitanic/titanic/services/stern/internal/server"
)

func PeppySkillIssue(ctx *server.Context) {
	ctx.RenderTemplate(
		http.StatusOK, "pages/public/peppy",
		buildDefaultView(ctx),
	)
}
