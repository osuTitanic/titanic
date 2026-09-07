package routes

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/osuTitanic/titanic/internal/activity"
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
	"github.com/redis/go-redis/v9"
)

const replayViewCooldown = 60 * time.Second

// /web/osu-getreplay.php -> View the replay of a score
func Replay(ctx *server.Context) {
	scoreId, err := ctx.QueryValueInt64("c")
	if err != nil {
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	score, err := ctx.State.Scores.ById(
		scoreId, "Beatmap.Beatmapset", "User",
	)
	if err != nil {
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}
	if score == nil {
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}
	if score.Hidden {
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}

	rawReplay, err := ctx.State.Storage.ReadStream(
		strconv.FormatInt(score.Id, 10),
		"replays",
	)
	if err != nil {
		ctx.Logger.Warn("Replay requested for score with no replay data", "score_id", score.Id)
		ctx.Response.WriteHeader(http.StatusNotFound)
		return
	}

	// If the user has authenticated themselves, we can increase the replay views
	// for this score and the score's player, as well as the event / activity
	user, _ := ctx.AuthenticateUserFromQuery(
		"u", "h", false,
	)
	if err = increaseReplayViews(user, score, ctx); err != nil {
		ctx.Logger.Warn(
			"Failed to increase replay views",
			"error", err, "score_id", score.Id, "viewer_id", user.Id,
		)
	}

	ctx.Response.Header().Set("Content-Type", "application/octet-stream")
	ctx.Response.WriteHeader(http.StatusOK)

	if _, err = io.Copy(ctx.Response, rawReplay); err != nil {
		ctx.Logger.Warn(
			"Failed to write replay data to response",
			"error", err, "score_id", score.Id,
		)
	}
}

func increaseReplayViews(viewer *schemas.User, score *schemas.Score, ctx *server.Context) error {
	if viewer == nil {
		return nil
	}
	if viewer.Id == score.User.Id {
		return nil
	}

	cooldownKey := fmt.Sprintf(
		"replay_cooldown:%d:%d",
		viewer.Id, score.User.Id,
	)
	cooldown, err := ctx.State.Redis.Get(
		ctx.Request.Context(),
		cooldownKey,
	).Result()

	if err != nil && err != redis.Nil {
		return err
	}
	if cooldown != "" {
		return nil
	}
	ctx.State.Redis.Set(
		ctx.Request.Context(),
		cooldownKey, "1", replayViewCooldown,
	)

	err = ctx.State.Histories.UpdateReplayViews(score.User.Id, score.Mode)
	if err != nil {
		return err
	}
	err = ctx.State.Stats.UpdateReplayViews(score.User.Id, score.Mode)
	if err != nil {
		return err
	}
	err = ctx.State.Scores.UpdateReplayViews(score.Id)
	if err != nil {
		return err
	}

	err = activity.Submit(
		ctx.State,
		viewer.Id,
		&score.Mode,
		constants.ActivityReplayWatched,
		map[string]any{
			"username":     viewer.Name,
			"score_id":     score.Id,
			"target_id":    score.User.Id,
			"target_name":  score.User.Name,
			"beatmap_id":   score.BeatmapId,
			"beatmap_name": score.Beatmap.Name(),
		},
		false, // should not be sent to #announce
		true,  // should be hidden in user profile
	)
	if err != nil {
		return err
	}
	return nil
}
