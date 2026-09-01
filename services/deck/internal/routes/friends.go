package routes

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

// /web/osu-getfriends.php -> Retrieve the authenticated user's friend IDs
func Friends(ctx *server.Context) {
	user, ok := ctx.HandleUserAuthenticationSimple("u", "h", false)
	if !ok {
		return
	}

	friendIds, err := ctx.State.Repositories.Relationships.TargetIdsByStatus(
		user.Id,
		constants.RelationshipStatusFriend,
	)
	if err != nil {
		ctx.Logger.Error("Failed to fetch friends", "user_id", user.Id, "error", err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: add some sort of shared concat / join function for integers
	var response strings.Builder
	for index, friendId := range friendIds {
		if index > 0 {
			response.WriteByte('\n')
		}
		response.WriteString(strconv.Itoa(friendId))
	}
	ctx.RenderText(http.StatusOK, response.String())
}
