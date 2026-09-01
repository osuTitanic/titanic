package routes

import (
	"net/http"
	"strings"

	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

// /web/osu-markasread.php -> Mark all direct messages from a user as read
func MarkAsRead(ctx *server.Context) {
	user, ok := ctx.HandleUserAuthenticationSimple("u", "h", false)
	if !ok {
		return
	}
	channel := ctx.QueryValue("channel")
	username := schemas.ResolveSafeName(strings.TrimSpace(channel))

	if strings.HasPrefix(channel, "#") {
		// We don't support marking public channels as read right now
		ctx.Response.WriteHeader(http.StatusOK)
		return
	}

	target, err := ctx.State.Users.BySafeName(username)
	if err != nil {
		ctx.Logger.Error("Failed to resolve direct message target", "channel", channel, "error", err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}
	if target == nil {
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}

	_, err = ctx.State.Messages.UpdatePrivateAll(
		&schemas.DirectMessage{
			SenderId: target.Id,
			TargetId: user.Id,
			Read:     true,
		},
		"read",
	)
	if err != nil {
		ctx.Logger.Error(
			"Failed to mark direct messages as read",
			"user_id", user.Id,
			"target_id", target.Id,
			"error", err,
		)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Logger.Info(
		"Direct messages marked as read",
		"user_id", user.Id,
		"target_id", target.Id,
	)
	ctx.Response.WriteHeader(http.StatusOK)
}
