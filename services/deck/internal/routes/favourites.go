package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/osuTitanic/titanic/internal/activity"
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

const (
	favouritesLimitMessage     = "You have too many favourite maps. Please go to your profile and delete some first."
	favouriteAlreadySetMessage = "You have already favourited this map..."
)

// /web/osu-getfavourites.php -> Get the user's favourite beatmapset IDs
func GetFavourites(ctx *server.Context) {
	user, ok := ctx.HandleUserAuthenticationSimple("u", "h", false)
	if !ok {
		return
	}
	user.LatestActivity = time.Now()
	ctx.State.Users.Update(user, "latest_activity")

	favourites, err := ctx.State.Favourites.ManyByUserId(user.Id)
	if err != nil {
		ctx.Logger.Error("Failed to fetch favourites", "user_id", user.Id, "error", err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response strings.Builder
	for index, favourite := range favourites {
		if index > 0 {
			response.WriteByte('\n')
		}
		response.WriteString(strconv.Itoa(favourite.SetId))
	}

	ctx.Logger.Info("Requested favourites", "user_id", user.Id, "count", len(favourites))
	ctx.RenderText(http.StatusOK, response.String())
}

// /web/osu-addfavourite.php -> Add a beatmapset to the user's favourites
func AddFavourite(ctx *server.Context) {
	setId, err := ctx.QueryValueInt("a")
	if err != nil {
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	user, ok := ctx.HandleUserAuthenticationSimple("u", "h", true)
	if !ok {
		return
	}
	user.LatestActivity = time.Now()
	ctx.State.Users.Update(user, "latest_activity")

	count, err := ctx.State.Favourites.CountByUserId(user.Id)
	if err != nil {
		ctx.Logger.Error("Failed to count favourites", "user_id", user.Id, "error", err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}
	if count >= ctx.State.Config.BeatmapFavoritesLimit {
		ctx.Logger.Warn("User has too many favourites", "user_id", user.Id, "count", count)
		ctx.RenderText(http.StatusOK, favouritesLimitMessage)
		return
	}

	beatmapset, err := ctx.State.Beatmapsets.ById(setId)
	if err != nil {
		ctx.Logger.Error("Failed to fetch beatmapset", "set_id", setId, "error", err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}
	if beatmapset == nil {
		ctx.Logger.Warn("Beatmapset was not found", "set_id", setId)
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}

	created, err := ctx.State.Favourites.CreateIfAbsent(&schemas.BeatmapFavourite{
		UserId: user.Id,
		SetId:  setId,
	})
	if err != nil {
		ctx.Logger.Error("Failed to create favourite", "user_id", user.Id, "set_id", setId, "error", err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !created {
		ctx.Logger.Warn("Beatmap was already favourited", "user_id", user.Id, "set_id", setId)
		ctx.RenderText(http.StatusOK, favouriteAlreadySetMessage)
		return
	}
	count++

	ctx.Logger.Info(
		"Added favourite",
		"user_id", user.Id, "set_id", setId, "count", count,
	)

	err = activity.Submit(
		ctx.State,
		user.Id,
		nil,
		constants.ActivityBeatmapFavouriteAdded,
		map[string]any{
			"username":        user.Name,
			"beatmapset_id":   beatmapset.Id,
			"beatmapset_name": fmt.Sprintf("%s (%s)", beatmapset.Name(), beatmapset.CreatorName()),
		},
		false, // should not be sent to #announce
		true,  // should be hidden in user profile
	)
	if err != nil {
		ctx.Logger.Warn("Failed to broadcast favourite activity", "user_id", user.Id, "set_id", setId, "error", err)
	}

	plural := ""
	if count > 1 {
		plural = "s"
	}
	ctx.RenderText(
		http.StatusOK,
		fmt.Sprintf("Added to favourites! You have a total of %d favourite%s", count, plural),
	)
}
