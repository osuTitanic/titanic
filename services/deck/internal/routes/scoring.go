package routes

import (
	"net/http"

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
	for _, warning := range result.Warnings {
		ctx.Logger.Warn("Score submission post-processing step failed", "error", warning)
	}

	if result.Rejected() {
		ctx.Logger.Debug("Score submission rejected", "result", result.Type)
		writeSubmissionError(ctx, endpoint, result.Type)
		return
	}
	writeSubmissionResponse(ctx, endpoint, result.Submission)
}

func writeSubmissionResponse(ctx *server.Context, endpoint scoring.Endpoint, result *scoring.SubmissionContext) {
	// TODO: Render score submission response
}

func writeSubmissionError(ctx *server.Context, endpoint scoring.Endpoint, resultType scoring.ResultType) {
	if resultType == scoring.ResultBanchoUnavailable {
		// Client will perform a delayed retry if an error occurs
		// Let's hope it will connect to bancho in the meantime
		ctx.Response.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	if !endpoint.UsesLegacyResponse() {
		// use `error: <type>` response
		ctx.RenderText(http.StatusOK, resultType.String())
		return
	}

	// /web/osu-submit.php has no custom error response strings
	// we'll just use http status codes

	status := http.StatusBadRequest
	switch resultType {
	case scoring.ResultUserNotFound,
		scoring.ResultInvalidPassword,
		scoring.ResultInactive,
		scoring.ResultBanned:
		status = http.StatusUnauthorized
	case scoring.ResultBeatmapUnavailable:
		status = http.StatusNotFound
	default:
		status = http.StatusBadRequest
	}
	ctx.Response.WriteHeader(status)
}
