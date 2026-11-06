package routes

import (
	"io"
	"strconv"
	"strings"

	"github.com/osuTitanic/titanic/services/stern/internal/server"
)

func BeatmapThumbnail(ctx *server.Context) {
	filename := ctx.PathValue("filename")

	// Handle filenames such as "1.jpg" (small) and "1l.jpg" (large)
	key, _, _ := strings.Cut(filename, ".")
	large := strings.Contains(key, "l")

	setId, err := strconv.Atoi(strings.ReplaceAll(key, "l", ""))
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}

	stream, err := ctx.State.Resources.Background(ctx.Context(), setId, large)
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}
	defer stream.Close()

	// If a cache key is provided, the thumbnail may be cached by the client
	if ctx.QueryValue("c") != "" {
		ctx.Response.Header().Set("Cache-Control", "public, max-age=604800, immutable")
	}

	ctx.Response.Header().Set("Content-Type", "image/jpeg")
	ctx.Response.WriteHeader(200)
	io.Copy(ctx.Response, stream)
}

func BeatmapAudioPreview(ctx *server.Context) {
	filename := ctx.PathValue("filename")

	// Handle filenames such as "1.mp3"
	key, _, _ := strings.Cut(filename, ".")

	setId, err := strconv.Atoi(key)
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}

	stream, err := ctx.State.Resources.Preview(ctx.Context(), setId)
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}
	defer stream.Close()

	// If a cache key is provided, the preview may be cached by the client
	if ctx.QueryValue("c") != "" {
		ctx.Response.Header().Set("Cache-Control", "public, max-age=604800, immutable")
	}

	ctx.Response.Header().Set("Content-Type", "audio/mpeg")
	ctx.Response.WriteHeader(200)
	io.Copy(ctx.Response, stream)
}
