package routes

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"

	"github.com/osuTitanic/titanic/internal/authentication"
	"github.com/osuTitanic/titanic/internal/repositories"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

const (
	beatmapInfoMaxMaps     = 100
	beatmapInfoMaxBodySize = 64 << 10
)

type beatmapInfoRequest struct {
	Filenames []string `json:"Filenames"`
	Ids       []int    `json:"Ids"`
}

// /web/osu-getbeatmapinfo.php -> Resolve beatmap metadata & pb's, similar to Bancho's BeatmapInfo packets
func BeatmapInfo(ctx *server.Context) {
	username := ctx.QueryValue("u")
	password := ctx.QueryValue("h")
	user, err := ctx.State.Repositories.Users.ByName(username)
	if err != nil {
		ctx.Logger.Error("Failed to fetch user", "username", username, "error", err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user == nil {
		ctx.Response.WriteHeader(http.StatusUnauthorized)
		return
	}
	if !authentication.VerifyPasswordHashFromMd5(password, user.Bcrypt) {
		ctx.Response.WriteHeader(http.StatusUnauthorized)
		return
	}

	request, ok := decodeBeatmapInfoRequest(ctx)
	if !ok {
		return
	}

	results, err := ctx.State.Repositories.Beatmaps.FetchInfoByFilenamesOrIds(
		user.Id,
		request.Filenames,
		request.Ids,
	)
	if err != nil {
		ctx.Logger.Error("Failed to fetch beatmap info", "error", err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx.RenderText(http.StatusOK, formatBeatmapInfo(request, results))
}

func formatBeatmapInfo(request beatmapInfoRequest, results []repositories.BeatmapInfoResult) string {
	var output strings.Builder

	for _, result := range results {
		if output.Len() > 0 {
			output.WriteByte('\n')
		}
		// The client will identify the beatmaps by their index
		// in the "beatmapInfoSendList" array for the filenames
		// For beatmap IDs, the client won't use the index, so -1 will suffice
		filenameIndex := slices.Index(request.Filenames, result.Filename)
		beatmapInfoReply := result.Format(filenameIndex)
		output.WriteString(beatmapInfoReply)
	}
	return output.String()
}

func decodeBeatmapInfoRequest(ctx *server.Context) (request beatmapInfoRequest, ok bool) {
	reader := http.MaxBytesReader(ctx.Response, ctx.Request.Body, beatmapInfoMaxBodySize)
	decoder := json.NewDecoder(reader)

	if err := decoder.Decode(&request); err != nil {
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return beatmapInfoRequest{}, false
	}
	if request.Filenames == nil || request.Ids == nil {
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return beatmapInfoRequest{}, false
	}

	totalMaps := len(request.Filenames) + len(request.Ids)
	if totalMaps > beatmapInfoMaxMaps {
		ctx.RenderText(http.StatusOK, "")
		return
	}
	return request, true
}
