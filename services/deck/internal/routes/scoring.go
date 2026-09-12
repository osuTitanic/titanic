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
		status, text := endpoint.FormatError(result.Type)
		ctx.RenderText(status, text)
		return
	}

	ctx.RenderText(http.StatusOK, endpoint.FormatResponse(result.Submission))
}
