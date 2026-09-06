package routes

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"

	"github.com/osuTitanic/titanic/internal/repositories"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

const (
	BeatmapInfoMaxMaps     = 100
	BeatmapInfoMaxBodySize = 64 << 10
)

type BeatmapInfoRequest struct {
	Filenames []string `json:"Filenames"`
	Ids       []int    `json:"Ids"`
}

func (r BeatmapInfoRequest) Total() int {
	return len(r.Filenames) + len(r.Ids)
}

func NewBeatmapInfoRequest(ctx *server.Context) (request BeatmapInfoRequest, ok bool) {
	reader := http.MaxBytesReader(ctx.Response, ctx.Request.Body, BeatmapInfoMaxBodySize)
	decoder := json.NewDecoder(reader)

	if err := decoder.Decode(&request); err != nil {
		return BeatmapInfoRequest{}, false
	}
	if request.Filenames == nil || request.Ids == nil {
		return BeatmapInfoRequest{}, false
	}
	if request.Total() > BeatmapInfoMaxMaps {
		return BeatmapInfoRequest{}, false
	}
	return request, true
}

// /web/osu-getbeatmapinfo.php -> Resolve beatmap metadata & pb's, similar to Bancho's BeatmapInfo packets
func BeatmapInfo(ctx *server.Context) {
	user, ok := ctx.HandleUserAuthenticationSimple("u", "h", true)
	if !ok {
		return
	}

	request, ok := NewBeatmapInfoRequest(ctx)
	if !ok {
		ctx.Response.WriteHeader(http.StatusBadRequest)
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

func formatBeatmapInfo(request BeatmapInfoRequest, results []repositories.BeatmapInfoResult) string {
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
