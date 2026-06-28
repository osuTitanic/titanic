package routes

import (
	"strconv"
	"strings"

	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
)

func Group(ctx *server.Context) {
	groupIdString := strings.TrimSpace(ctx.PathValue("id"))
	groupId, err := strconv.Atoi(groupIdString)
	if err != nil {
		NotFound(ctx)
		return
	}

	group, err := ctx.State.Groups.ById(groupId)
	if err != nil {
		InternalServerError(ctx)
		return
	}
	if group == nil {
		NotFound(ctx)
		return
	}

	view := templates.GroupView{
		DefaultView: buildDefaultView(ctx),
		Group:       group,
	}
	ctx.RenderTemplate(200, "pages/public/group", view)
}
