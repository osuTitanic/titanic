package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/repositories"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

// /web/osu-getstatus.php -> Resolve beatmap submission statuses by checksum
func BeatmapStatus(ctx *server.Context) {
	if !ctx.Request.URL.Query().Has("c") {
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	checksums := strings.Split(ctx.QueryValue("c"), ",")
	if len(checksums) > 60 {
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, checksum := range checksums {
		if checksum != "" && len(checksum) != 32 {
			ctx.Response.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	results, err := ctx.State.Repositories.Beatmaps.FetchStatusesByChecksums(checksums)
	if err != nil {
		ctx.Logger.Error("Failed to fetch beatmap statuses", "error", err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.RenderText(http.StatusOK, formatBeatmapStatuses(checksums, results))
}

func formatBeatmapStatuses(checksums []string, results []repositories.BeatmapStatusResult) string {
	resultsByChecksum := make(map[string]repositories.BeatmapStatusResult, len(results))
	for _, result := range results {
		resultsByChecksum[result.Checksum] = result
	}

	response := make([]string, 0, len(results))
	for _, checksum := range checksums {
		beatmap, exists := resultsByChecksum[checksum]
		if !exists {
			continue
		}

		ranked := 0
		if beatmap.Status > constants.BeatmapStatusPending {
			ranked = 1
		}

		topicId := ""
		if beatmap.TopicId != nil {
			topicId = strconv.Itoa(*beatmap.TopicId)
		}
		response = append(response, fmt.Sprintf(
			"%s,%d,%d,%d,%s",
			checksum,
			ranked,
			beatmap.Id,
			beatmap.SetId,
			topicId,
		))
	}
	return strings.Join(response, "\n")
}
