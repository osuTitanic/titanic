package routes

import (
	"io"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

// /web/maps/{query} -> Download a beatmap file by id, filename or checksum
func BeatmapFile(ctx *server.Context) {
	query := strings.TrimSpace(ctx.PathValue("query"))
	beatmap, err := resolveBeatmap(query, ctx)
	if err != nil {
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}
	if beatmap == nil {
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}

	stream, err := ctx.State.Resources.Osu(beatmap.Id)
	if err != nil {
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}
	defer stream.Close()

	ctx.Response.Header().Set(
		"Content-Type",
		"application/octet-stream",
	)
	ctx.Response.Header().Set(
		"Content-Disposition",
		mime.FormatMediaType("attachment", map[string]string{"filename": beatmap.Filename}),
	)
	ctx.Response.Header().Set(
		"Last-Modified",
		beatmap.LastUpdate.UTC().Format(http.TimeFormat),
	)
	ctx.Response.WriteHeader(http.StatusOK)
	io.Copy(ctx.Response, stream)
}

func resolveBeatmap(query string, ctx *server.Context) (*schemas.Beatmap, error) {
	if id, err := strconv.Atoi(query); err == nil {
		return ctx.State.Repositories.Beatmaps.ById(id)
	}
	if strings.HasSuffix(query, ".osu") {
		return ctx.State.Repositories.Beatmaps.ByFilename(query)
	}
	return ctx.State.Repositories.Beatmaps.ByChecksum(query)
}
