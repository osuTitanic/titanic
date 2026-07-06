package routes

import (
	"net/http"

	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
)

func AccountOverview(ctx *server.Context) {
	if !ctx.RequireLogin() {
		return
	}
	userId := ctx.CurrentUser.Id

	totalPosts, err := ctx.State.Repositories.ForumPosts.CountByUserId(userId)
	if err != nil {
		ctx.Logger.Error("Failed to fetch post count", "user", userId, "error", err)
		InternalServerError(ctx)
		return
	}

	notifications, err := ctx.State.Repositories.Notifications.UnreadByUserId(userId)
	if err != nil {
		ctx.Logger.Error("Failed to fetch notifications", "user", userId, "error", err)
		InternalServerError(ctx)
		return
	}

	bookmarks, err := ctx.State.Repositories.ForumBookmarks.FetchByUserId(userId, "Topic", "Topic.Forum")
	if err != nil {
		ctx.Logger.Error("Failed to fetch bookmarks", "user", userId, "error", err)
		InternalServerError(ctx)
		return
	}

	view := templates.SettingsOverviewView{
		DefaultView:   buildDefaultView(ctx),
		TotalPosts:    totalPosts,
		Notifications: notifications,
		Bookmarks:     bookmarks,
	}
	ctx.RenderTemplate(http.StatusOK, "pages/account/settings_overview", view)
}
