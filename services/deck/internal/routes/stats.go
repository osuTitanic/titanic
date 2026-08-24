package routes

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"time"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

// /web/osu-stat.php -> Retrieve own user stats, used by IRC osu! clients
func UserStats(ctx *server.Context) {
	username := ctx.QueryValue("u")
	password := ctx.QueryValue("p")

	if password == "" {
		// NOTE: We don't really bother to validate the password of the user here
		ctx.Response.WriteHeader(http.StatusUnauthorized)
		return
	}

	renderStatsResponse(username, ctx)
}

// /web/osu-statoth.php -> Retrieve user stats of another player, used by IRC osu! clients
func UserStatsOther(ctx *server.Context) {
	username := ctx.QueryValue("u")
	checksum := ctx.QueryValue("c")

	if checksum == "" {
		ctx.Response.WriteHeader(http.StatusForbidden)
		return
	}

	checksumData := fmt.Sprintf("%sprettyplease!!!", username)
	checksumMatch := md5.Sum([]byte(checksumData))

	if fmt.Sprintf("%x", checksumMatch) != checksum {
		ctx.Response.WriteHeader(http.StatusForbidden)
		return
	}

	renderStatsResponse(username, ctx)
}

func renderStatsResponse(username string, ctx *server.Context) {
	userId, err := ctx.State.Repositories.Users.GetUserId(username)
	if err != nil {
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}

	avatarChecksum, err := resolveAvatarChecksum(userId, ctx)
	if err != nil {
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}

	currentRank, err := ctx.State.Rankings.GlobalRank(userId, constants.ModeOsu)
	if err != nil {
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}

	currentAcc, err := ctx.State.Rankings.Accuracy(userId, constants.ModeOsu)
	if err != nil {
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}

	currentScore, err := ctx.State.Rankings.Score(userId, constants.ModeOsu)
	if err != nil {
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Cap score to a signed 32-bit integer to prevent overflow
	currentScoreCapped := min(currentScore, 2147483647)

	response := fmt.Sprintf(
		"%d|%.8f|%d|%d|%d|%d_%s.png",
		currentScoreCapped,
		currentAcc,
		currentScore, // NOTE: This field is usually empty & unused
		userId,       //	   Same goes for this field
		currentRank,
		userId, avatarChecksum, // Avatar Filename
	)
	ctx.RenderText(http.StatusOK, response)
}

func resolveAvatarChecksum(userId int, ctx *server.Context) (string, error) {
	// Check if the avatar checksum is already cached
	result := ctx.State.Redis.Get(
		ctx.Request.Context(),
		fmt.Sprintf("bancho:avatar_hash:%d", userId),
	)

	if result.Err() == nil {
		return result.Val(), nil
	}

	// Fallback to database query & cache the result if successful
	checksum, err := ctx.State.Repositories.Users.GetAvatarChecksum(userId)
	if err != nil {
		return "", err
	}

	ctx.State.Redis.Set(
		ctx.Request.Context(),
		fmt.Sprintf("bancho:avatar_hash:%d", userId),
		checksum, time.Hour*24,
	)
	return checksum, nil
}
