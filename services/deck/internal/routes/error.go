package routes

import (
	"net/http"

	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

// /web/osu-error.php -> osu! error report handler
func ErrorReport(ctx *server.Context) {
	// Client may submit the form fields:
	// - bc (beatmap md5)
	// - b (beatmap id)
	// - u (username)
	// - i (user id)
	// - culture
	// - gamemode
	// - gametime
	// - config
	// - iltrace
	// - exehash
	// - feedback
	// - stacktrace
	// - audiotime
	// - exception
	// - version

	// We don't really do anything with these, since
	// we aren't going to fix these anyways :)
	ctx.Response.WriteHeader(http.StatusOK)
}
