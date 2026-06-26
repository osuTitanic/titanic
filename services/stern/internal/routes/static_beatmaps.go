package routes

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
)

func BeatmapThumbnail(ctx *server.Context) {
	filename := ctx.PathValue("filename")

	// Handle filenames such as "1.jpg" (small) and "1l.jpg" (large)
	key := strings.SplitN(filename, ".", 2)[0]
	large := strings.Contains(key, "l")

	setId, err := strconv.Atoi(strings.ReplaceAll(key, "l", ""))
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}

	stream, err := ctx.State.Resources.Background(setId, large)
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}
	defer stream.Close()

	// If a cache key is provided, the thumbnail may be cached by the client
	if ctx.Request.URL.Query().Get("c") != "" {
		ctx.Response.Header().Set("Cache-Control", "public, max-age=86400")
	}

	ctx.Response.Header().Set("Content-Type", "image/jpeg")
	ctx.Response.WriteHeader(200)
	io.Copy(ctx.Response, stream)
}

func BeatmapAudioPreview(ctx *server.Context) {
	filename := ctx.PathValue("filename")

	// Handle filenames such as "1.mp3"
	key := strings.SplitN(filename, ".", 2)[0]

	setId, err := strconv.Atoi(key)
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}

	stream, err := ctx.State.Resources.Preview(setId)
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}
	defer stream.Close()

	// If a cache key is provided, the preview may be cached by the client
	if ctx.Request.URL.Query().Get("c") != "" {
		ctx.Response.Header().Set("Cache-Control", "public, max-age=86400")
	}

	ctx.Response.Header().Set("Content-Type", "audio/mpeg")
	ctx.Response.WriteHeader(200)
	io.Copy(ctx.Response, stream)
}

func BeatmapDownload(ctx *server.Context) {
	filename := ctx.PathValue("filename")

	// Handle filenames such as "1 Kenji Ninuma - DISCO PRINCE.osz"
	filename = strings.SplitN(filename, " ", 2)[0]
	noVideo := strings.Contains(filename, "n")

	setIdString := strings.TrimSuffix(filename, "n")
	setId, err := strconv.Atoi(setIdString)
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}

	beatmapset, err := ctx.State.Beatmapsets.ById(setId)
	if err != nil || beatmapset == nil {
		ctx.Response.WriteHeader(404)
		return
	}
	if !beatmapset.Available {
		ctx.Response.WriteHeader(451)
		return
	}

	// noVideo can only be true if the beatmapset has videos
	noVideo = noVideo && beatmapset.HasVideo

	oszStream, oszSize, err := ctx.State.Resources.Osz(setId, noVideo)
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}

	oszFilename := fmt.Sprintf("%d %s.osz", beatmapset.Id, beatmapset.Name())
	if noVideo {
		oszFilename = fmt.Sprintf("%d %s (no video).osz", beatmapset.Id, beatmapset.Name())
	}

	ctx.Response.Header().Set("Content-Type", "application/octet-stream")
	ctx.Response.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", oszFilename))
	ctx.Response.Header().Set("Last-Modified", beatmapset.LastUpdate.Format("Mon, 02 Jan 2006 15:04:05 GMT"))

	if oszSize > 0 {
		// Set content length if we can determine it, otherwise we'll use chunked transfer encoding
		ctx.Response.Header().Set("Content-Length", strconv.FormatInt(oszSize, 10))
	}

	ctx.Response.WriteHeader(200)
	io.Copy(ctx.Response, oszStream)
}
