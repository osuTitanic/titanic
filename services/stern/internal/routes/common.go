package routes

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/osuTitanic/titanic-go/internal/state"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
)

// RenderError renders a generic message page for errors
func RenderError(ctx *server.Context, status int, heading, message string) {
	view := templates.ErrorMessageView{
		DefaultView: buildDefaultView(ctx),
		Title:       fmt.Sprintf("%s - Titanic!", heading),
		Heading:     heading,
		Message:     message,
	}
	ctx.RenderTemplate(status, "errors/generic", view)
}

// RenderErrorPage renders a custom error page based on the provided template name
func RenderErrorPage(ctx *server.Context, status int, name string) {
	ctx.RenderTemplate(
		status,
		"errors/custom/"+name,
		buildDefaultView(ctx),
	)
}

func NotFound(ctx *server.Context) {
	ctx.RenderTemplate(
		http.StatusNotFound,
		"errors/404",
		buildDefaultView(ctx),
	)
}

func InternalServerError(ctx *server.Context) {
	if ctx.Templates == nil {
		ctx.Logger.Error("Failed to render template", "template", "errors/500", "error", "templates engine is not configured")
		templates.InternalServerErrorFallback(ctx.Response)
		return
	}

	body, err := ctx.Templates.Render("errors/500", buildDefaultView(ctx))
	if err != nil {
		ctx.Logger.Error("Failed to render template", "template", "errors/500", "error", err)
		templates.InternalServerErrorFallback(ctx.Response)
		return
	}

	ctx.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.Response.WriteHeader(http.StatusInternalServerError)
	if _, err := ctx.Response.Write(body); err != nil {
		ctx.Logger.Error("Failed to write response body", "template", "errors/500", "error", err)
	}
}

func BeatmapNotFound(ctx *server.Context) {
	RenderErrorPage(ctx, http.StatusNotFound, "beatmap_not_found")
}

func ForumNotFound(ctx *server.Context) {
	RenderErrorPage(ctx, http.StatusNotFound, "forum_not_found")
}

func TopicNotFound(ctx *server.Context) {
	RenderErrorPage(ctx, http.StatusNotFound, "topic_not_found")
}

func PostNotFound(ctx *server.Context) {
	RenderErrorPage(ctx, http.StatusNotFound, "post_not_found")
}

func UserNotFound(ctx *server.Context) {
	RenderErrorPage(ctx, http.StatusNotFound, "user_not_found")
}

func TopicLocked(ctx *server.Context) {
	RenderErrorPage(ctx, http.StatusForbidden, "topic_locked")
}

func PostLocked(ctx *server.Context) {
	RenderErrorPage(ctx, http.StatusForbidden, "post_locked")
}

func UserSilenced(ctx *server.Context) {
	RenderErrorPage(ctx, http.StatusForbidden, "user_silenced")
}

func UserRestricted(ctx *server.Context) {
	RenderErrorPage(ctx, http.StatusForbidden, "user_restricted")
}

func PostingTooQuickly(ctx *server.Context) {
	RenderErrorPage(ctx, http.StatusTooManyRequests, "posting_too_quickly")
}

func PostTooLong(ctx *server.Context) {
	RenderErrorPage(ctx, http.StatusBadRequest, "post_too_long")
}

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
