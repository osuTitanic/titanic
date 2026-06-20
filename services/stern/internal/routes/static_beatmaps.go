package routes

import (
	"io"
	"strings"

	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
)

func BeatmapThumbnail(ctx *server.Context) {
	filename := ctx.PathValue("filename")
	key := strings.SplitN(filename, ".", 2)[0]

	avatar, err := ctx.State.Storage.ReadStream(key, "thumbnails")
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}

	ctx.Response.Header().Set("Content-Type", "image/jpeg")
	ctx.Response.WriteHeader(200)
	io.Copy(ctx.Response, avatar)
}

func BeatmapAudioPreview(ctx *server.Context) {
	filename := ctx.PathValue("filename")
	key := strings.SplitN(filename, ".", 2)[0]

	avatar, err := ctx.State.Storage.ReadStream(key, "audio")
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}

	ctx.Response.Header().Set("Content-Type", "audio/mpeg")
	ctx.Response.WriteHeader(200)
	io.Copy(ctx.Response, avatar)
}
