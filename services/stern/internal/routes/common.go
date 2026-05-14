package routes

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/osuTitanic/titanic-go/internal/state"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
)

func NotFound(ctx *server.Context) {
	ctx.RenderTemplate(
		http.StatusNotFound, "errors/404",
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

func buildDefaultView(ctx *server.Context) templates.DefaultView {
	currentPath := ctx.Request.URL.Path
	currentURI := "/"

	if requestURI := ctx.Request.URL.RequestURI(); requestURI != "" {
		currentURI = requestURI
	}

	return templates.DefaultView{
		Stats:       buildStatistics(ctx.State),
		Query:       ctx.Request.URL.Query(),
		Config:      ctx.State.Config,
		CSRFToken:   ctx.CSRFToken,
		CurrentUser: ctx.CurrentUser,
		CurrentPath: currentPath,
		CurrentURI:  currentURI,
	}
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
