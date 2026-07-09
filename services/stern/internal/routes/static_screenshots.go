package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/osuTitanic/titanic/services/stern/internal/server"
)

func Screenshot(ctx *server.Context) {
	id, err := ctx.PathValueInt("id")
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}

	screenshot, err := ctx.State.Screenshots.ById(id)
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}
	if screenshot.Hidden {
		ctx.Response.WriteHeader(404)
		return
	}

	if ctx.PathValue("checksum") != screenshot.Checksum() {
		ctx.Response.WriteHeader(404)
		return
	}

	image, err := ctx.State.Storage.Read(strconv.Itoa(id), "screenshots")
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}

	contentType := http.DetectContentType(image)
	fileExtension := strings.TrimPrefix(contentType, "image/")

	if !strings.HasPrefix(contentType, "image/") {
		ctx.Response.WriteHeader(404)
		return
	}

	filename := fmt.Sprintf(
		"ss (%s).%s", // inspired by puush
		screenshot.CreatedAt.Format("2006-01-02 at 15.04.05"), fileExtension,
	)

	ctx.Response.Header().Set("Content-Type", contentType)
	ctx.Response.Header().Set("Cache-Control", "public, max-age=1209600, immutable")
	ctx.Response.Header().Set("Date", screenshot.CreatedAt.Format("Mon, 02 Jan 2006 15:04:05 GMT"))
	ctx.Response.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))
	ctx.Response.Header().Set("Content-Length", strconv.Itoa(len(image)))
	ctx.Response.WriteHeader(200)
	ctx.Response.Write(image)
}

func ScreenshotRedirect(ctx *server.Context) {
	id, err := ctx.PathValueInt("id")
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}

	screenshot, err := ctx.State.Screenshots.ById(id)
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}
	if screenshot.Hidden {
		ctx.Response.WriteHeader(404)
		return
	}

	if time.Since(screenshot.CreatedAt) > 7*24*time.Hour {
		ctx.Response.WriteHeader(404)
		return
	}

	ctx.Redirect(301, fmt.Sprintf("/ss/%d/%s", id, screenshot.Checksum()))
}
