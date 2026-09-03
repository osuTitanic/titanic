package routes

import (
	"errors"
	"net/http"
	"time"

	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

const loginLockTTL = 10 * time.Second

// /web/osu-login.php -> Authenticate IRC-based osu! clients before they connect
//
// Before osu! had it's own game server (Bancho), it used an IRC server to handle in-game chatting.
// When a user logs in for the first time, it authenticates through this endpoint.
// After the first login was successful, it does not call this endpoint again.
// (this is also kinda annoying since the clients don't send any authentication info to the server)
func LegacyLogin(ctx *server.Context) {
	user, err := ctx.AuthenticateUserFromQuery("username", "password", false)
	switch {
	case errors.Is(err, server.ErrUserNotFound), errors.Is(err, server.ErrInvalidPassword):
		ctx.RenderText(http.StatusOK, "0")
		return
	case err != nil:
		ctx.Logger.Error("Failed to authenticate user", "error", err)
		ctx.RenderText(http.StatusInternalServerError, "0")
		return
	}

	if user.Restricted {
		ctx.RenderText(http.StatusOK, "0")
		return
	}
	if !user.Activated {
		ctx.RenderText(http.StatusOK, "0")
		return
	}
	user.LatestActivity = time.Now()
	ctx.State.Users.Update(user, "latest_activity")

	// On Titanic, we usually ask for an IRC token upon login
	// However, if the client has called this endpoint in the right
	// time frame we allow them to connect without a token using this lock
	key := "bancho:irc_login:" + user.SafeName
	if err := ctx.State.Redis.Set(ctx.Request.Context(), key, ctx.IP(), loginLockTTL).Err(); err != nil {
		ctx.Logger.Error("Failed to create login lock", "user_id", user.Id, "error", err)
		ctx.RenderText(http.StatusInternalServerError, "0")
		return
	}

	ctx.Logger.Info("User is connecting through IRC-based osu!", "user_id", user.Id)
	ctx.RenderText(http.StatusOK, "1")
}
