package routes

import (
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
)

func Group(ctx *server.Context) {
	groupId, err := ctx.PathValueInt("id")
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
	if group.Hidden {
		NotFound(ctx)
		return
	}

	groupUsers, err := ctx.State.Users.ManyByGroupId(group.Id)
	if err != nil {
		InternalServerError(ctx)
		return
	}

	view := templates.GroupView{
		DefaultView: buildDefaultView(ctx),
		Group:       group,
		Users:       groupUsers,
	}
	ctx.RenderTemplate(200, "pages/public/group", view)
}
