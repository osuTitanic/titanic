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

func BeatmapDownload(ctx *server.Context) {
	// Handle filenames such as "1 Kenji Ninuma - DISCO PRINCE.osz"
	filename := ctx.PathValue("filename")
	filename = strings.SplitN(filename, " ", 2)[0]
	noVideo := strings.Contains(filename, "n")

	setIdString := strings.TrimSuffix(filename, "n")
	setId, err := strconv.Atoi(setIdString)
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}

	beatmapset, err := ctx.State.Beatmapsets.ById(setId)
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}
	if !beatmapset.Available {
		ctx.Response.WriteHeader(451)
		return
	}

	// TODO: Retrieve osz from beatmap resource api
	oszStream, err := ctx.State.Storage.ReadStream(filename, "osz")
	if err != nil {
		ctx.Response.WriteHeader(404)
		return
	}

	// noVideo can only be true if the beatmapset has videos
	noVideo = noVideo && beatmapset.HasVideo

	oszFilename := fmt.Sprintf("%d %s.osz", beatmapset.Id, beatmapset.Name())
	oszSizeEstimated := beatmapset.OszFilesize

	if noVideo {
		oszFilename = fmt.Sprintf("%d %s (no video).osz", beatmapset.Id, beatmapset.Name())
		oszSizeEstimated -= beatmapset.OszFilesizeNovideo
	}
	// TODO: Populate osz sizes

	ctx.Response.Header().Set("Content-Type", "application/octet-stream")
	ctx.Response.Header().Set("Content-Length", strconv.Itoa(oszSizeEstimated))
	ctx.Response.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", oszFilename))
	ctx.Response.Header().Set("Last-Modified", beatmapset.LastUpdate.Format("Mon, 02 Jan 2006 15:04:05 GMT"))
	ctx.Response.WriteHeader(200)
	io.Copy(ctx.Response, oszStream)
}
