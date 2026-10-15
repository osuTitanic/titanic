package routes

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/osuTitanic/titanic/internal/discord"
	"github.com/osuTitanic/titanic/services/deck/internal/scoring"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

func SubmitScore(ctx *server.Context) {
	submitScore(ctx, scoring.EndpointLegacy)
}

func SubmitScoreModular(ctx *server.Context) {
	submitScore(ctx, scoring.EndpointModular)
}

func SubmitScoreModularSelector(ctx *server.Context) {
	submitScore(ctx, scoring.EndpointModularSelector)
}

func submitScore(ctx *server.Context, endpoint scoring.Endpoint) {
	submission, err := scoring.ResolveSubmissionContext(ctx, endpoint)
	if err != nil {
		ctx.Logger.Warn("Failed to parse score submission", "error", err)
		ctx.RenderText(http.StatusBadRequest, scoring.ResultRejected.String())
		return
	}
	password := endpoint.RequestValue(ctx.Request, "pass")

	processor := scoring.NewProcessor(ctx, submission)
	result, err := processor.Process(password)
	if err != nil {
		ctx.Logger.Error("Failed to process score submission", "error", err)
		ctx.RenderText(http.StatusInternalServerError, "")
		return
	}

	// Log all score submission warnings in the
	// console & send it to officer for staff
	defer logScoreWarnings(ctx, result)

	// After writing the response, we want to do some additional steps
	// like uploading the replay, broadcasting announcements, etc.
	defer processor.PostProcess()

	status := http.StatusOK
	response := endpoint.FormatResponse(result.Submission)

	if result.Rejected() {
		ctx.Logger.Debug("Score submission rejected", "result", result.Type)
		status, response = endpoint.FormatError(result.Type)
	}

	if err := writeSubmissionResponse(ctx, status, response); err != nil {
		ctx.Logger.Warn("Failed to write score submission response", "error", err)
	}
}

func writeSubmissionResponse(ctx *server.Context, status int, response string) error {
	ctx.Response.Header().Set("Content-Length", strconv.Itoa(len(response)))
	writeErr := ctx.RenderText(status, response)
	flushErr := http.NewResponseController(ctx.Response).Flush()
	return errors.Join(writeErr, flushErr)
}

func logScoreWarnings(ctx *server.Context, result scoring.Result) {
	if len(result.Warnings) == 0 {
		return
	}
	for _, warning := range result.Warnings {
		ctx.Logger.Warn("Score submission warning", "error", warning)
	}

	submission := result.Submission
	if submission == nil {
		return
	}
	score := submission.Score
	if score == nil {
		return
	}

	serverUrl := strings.TrimRight(ctx.State.Config.OsuBaseUrl(), "/")
	playerName := submission.Username
	beatmapName := submission.BeatmapChecksum
	clientName := score.ClientString

	if score.User != nil {
		playerName = score.User.Name
	}
	if score.Beatmap != nil {
		beatmapName = score.Beatmap.Name()
	}
	if clientName == "" {
		clientName = fmt.Sprintf("b%d", score.ClientVersion)
	}

	outcome := "Accepted"
	if result.Rejected() {
		outcome = fmt.Sprintf("Rejected (%s)", result.Type.String())
	}

	embed := discord.Embed{
		Title:     new("Score submission"),
		Color:     new(0xF2A900),
		Timestamp: new(time.Now().UTC()),
	}
	embed.AddField("Outcome", outcome, true)
	embed.AddField("Player", fmt.Sprintf("[%s](https://%s/u/%d)", playerName, serverUrl, score.UserId), true)
	embed.AddField("Beatmap", fmt.Sprintf("[%s](https://%s/u/%d)", beatmapName, serverUrl, score.BeatmapId), true)
	embed.AddField("Client", fmt.Sprintf("`%s`", clientName), true)

	if score.Id > 0 {
		embed.URL = new(fmt.Sprintf("%s/scores/%d", serverUrl, score.Id))
		embed.AddField("Score ID", fmt.Sprintf("[`%d`](%s)", score.Id, *embed.URL), true)
	}

	embeds := make([]discord.Embed, 1, len(result.Warnings)+1)
	embeds[0] = embed
	for _, warning := range result.Warnings {
		embeds = append(embeds, discord.Embed{
			Title:       new("Warning"),
			Description: new(warning.Error()),
			Color:       new(0xF2A900),
		})
	}
	notificationContext := context.WithoutCancel(ctx.Context())

	err := ctx.State.Officer.CallContext(
		notificationContext,
		discord.OfficerTagScoreSubmission,
		"",
		embeds...,
	)
	if err != nil {
		ctx.Logger.Warn(
			"Failed to send score submission officer notification",
			"user_id", score.UserId, "score_id", score.Id, "error", err,
		)
	}
	// TODO: Replay uploads?
}
