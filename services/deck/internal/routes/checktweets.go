package routes

import (
	"errors"
	"net/http"

	"github.com/osuTitanic/titanic/services/deck/internal/server"
	"github.com/redis/go-redis/v9"
)

// /web/osu-checktweets.php -> Bancho status message
//
// This endpoint was used to fetch tweets from @osustatus on
// twitter, which would be displayed on the client side.
func CheckTweets(ctx *server.Context) {
	message, err := ctx.State.Redis.Get(
		ctx.Request.Context(),
		"bancho:statusmessage",
	).Result()

	switch {
	case errors.Is(err, redis.Nil):
		ctx.RenderText(http.StatusOK, "")
	case err != nil:
		ctx.Logger.Error("Failed to retrieve bancho status message", "error", err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
	default:
		ctx.RenderText(http.StatusOK, message)
	}
}
