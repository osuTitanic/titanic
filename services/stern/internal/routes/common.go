package routes

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"github.com/osuTitanic/titanic/internal/state"
	"github.com/osuTitanic/titanic/services/stern/internal/server"
	"github.com/osuTitanic/titanic/services/stern/internal/templates"
)

func buildDefaultView(ctx *server.Context) templates.DefaultView {
	currentPath := ctx.Request.URL.Path
	currentURI := "/"

	if requestURI := ctx.Request.URL.RequestURI(); requestURI != "" {
		currentURI = requestURI
	}
	userAgent := ctx.Request.UserAgent()

	return templates.DefaultView{
		Stats:             buildStatistics(ctx.State),
		Query:             ctx.Request.URL.Query(),
		Config:            ctx.State.Config,
		CSRFToken:         ctx.CSRFToken,
		CurrentUser:       ctx.CurrentUser,
		CurrentPath:       currentPath,
		CurrentURI:        currentURI,
		NotificationCount: fetchNotificationCount(ctx),
		IsModernBrowser:   server.IsModernBrowser(userAgent),
		IsIE:              server.IsInternetExplorer(userAgent),
	}
}

func fetchNotificationCount(ctx *server.Context) int {
	if ctx.CurrentUser == nil {
		return 0
	}

	count, err := ctx.State.Notifications.CountUnreadByUserId(ctx.CurrentUser.Id)
	if err != nil {
		ctx.Logger.Error("Failed to fetch notification count", "error", err)
		return 0
	}
	return count
}

func buildDefaultViewWithPermissions(ctx *server.Context) templates.DefaultView {
	view := buildDefaultView(ctx)
	view.Permissions = ctx.Permissions()
	return view
}

func buildStatistics(state *state.State) (stats *templates.Statistics) {
	stats = &templates.Statistics{
		TotalUsers:     0,
		TotalScores:    0,
		OnlineUsersOsu: 0,
		OnlineUsersIrc: 0,
	}

	values, err := state.Redis.MGet(context.TODO(),
		"bancho:totalusers",
		"bancho:totalscores",
		"bancho:activity:osu",
		"bancho:activity:irc",
	).Result()
	if err != nil {
		state.Logger.Error("Failed to fetch statistics from redis", "error", err)
		return stats
	}

	if totalUsers, ok := values[0].(string); ok {
		stats.TotalUsers, _ = strconv.Atoi(totalUsers)
	}
	if totalScores, ok := values[1].(string); ok {
		stats.TotalScores, _ = strconv.Atoi(totalScores)
	}
	if onlineUsers, ok := values[2].(string); ok {
		stats.OnlineUsersOsu, _ = strconv.Atoi(onlineUsers)
	}
	if onlineUsers, ok := values[3].(string); ok {
		stats.OnlineUsersIrc, _ = strconv.Atoi(onlineUsers)
	}
	return stats
}

func sanitizeRedirectTarget(target string) string {
	target = strings.TrimSpace(target)
	if target == "" {
		return ""
	}

	parsed, err := url.Parse(target)
	if err != nil || parsed.IsAbs() || parsed.Host != "" {
		return "/"
	}

	requestURI := parsed.RequestURI()
	if requestURI == "" || !strings.HasPrefix(requestURI, "/") || strings.HasPrefix(requestURI, "//") {
		return "/"
	}

	return requestURI
}
