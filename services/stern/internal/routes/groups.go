package routes

import (
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/services/stern/internal/server"
	"github.com/osuTitanic/titanic/services/stern/internal/templates"
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

	users, err := ctx.State.Users.ManyByGroupId(group.Id)
	if err != nil {
		InternalServerError(ctx)
		return
	}

	// Exclude deactivated / restricted accounts
	visibleUsers := make([]*schemas.User, 0, len(users))
	for _, user := range users {
		if canViewUser(ctx, user) {
			visibleUsers = append(visibleUsers, user)
		}
	}

	view := templates.GroupView{
		DefaultView: buildDefaultView(ctx),
		Group:       group,
		Users:       visibleUsers,
	}
	ctx.RenderTemplate(200, "pages/public/group", view)
}
