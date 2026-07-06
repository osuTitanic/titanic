package routes

import (
	"fmt"
	"net/http"

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
	// We want to have as many failsafes as possible here, since
	// this is the last resort for rendering an error page.
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
