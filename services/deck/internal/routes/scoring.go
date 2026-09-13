package routes

import (
	"errors"
	"net/http"
	"strconv"

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
		ctx.Logger.Warn("Score submission warning", "error", warning)
		// TODO: Officer logging call
	}

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
