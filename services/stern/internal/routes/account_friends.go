package routes

import (
	"net/http"
	"slices"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/services/stern/internal/server"
	"github.com/osuTitanic/titanic/services/stern/internal/templates"
)

func AccountFriends(ctx *server.Context) {
	if !ctx.RequireLogin() {
		return
	}
	userId := ctx.CurrentUser.Id

	friends, err := ctx.State.Repositories.Relationships.FetchTargetUsers(userId, constants.RelationshipStatusFriend)
	if err != nil {
		ctx.Logger.Error("Failed to fetch friends", "user", userId, "error", err)
		InternalServerError(ctx)
		return
	}

	foes, err := ctx.State.Repositories.Relationships.FetchTargetUsers(userId, constants.RelationshipStatusFoe)
	if err != nil {
		ctx.Logger.Error("Failed to fetch foes", "user", userId, "error", err)
		InternalServerError(ctx)
		return
	}

	// Users that have added the current user back, used to flag mutual friends
	addedBy, err := ctx.State.Repositories.Relationships.UserIdsByStatus(userId, constants.RelationshipStatusFriend)
	if err != nil {
		ctx.Logger.Error("Failed to fetch mutual friends", "user", userId, "error", err)
		InternalServerError(ctx)
		return
	}

	entries := make([]*templates.SettingsFriend, 0, len(friends))
	for _, friend := range friends {
		mutual := slices.Contains(addedBy, friend.Id) ||
			slices.Contains(ctx.State.Config.SuperFriendlyUsers, friend.Id)
		entries = append(entries, &templates.SettingsFriend{
			User:     friend,
			IsMutual: mutual,
		})
	}

	view := templates.SettingsFriendsView{
		DefaultView: buildDefaultView(ctx),
		Friends:     entries,
		Foes:        foes,
	}
	ctx.RenderTemplate(http.StatusOK, "pages/account/settings_friends", view)
}
