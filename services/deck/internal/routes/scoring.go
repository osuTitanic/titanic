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

	// TODO: Process the submission
	_ = submission
	_ = password
	ctx.Response.WriteHeader(http.StatusNotImplemented)
}
