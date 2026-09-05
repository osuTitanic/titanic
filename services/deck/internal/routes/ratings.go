package routes

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/osuTitanic/titanic/internal/activity"
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

type RatingResultType uint8

const (
	RatingAuthenticationFailed RatingResultType = iota
	RatingBeatmapNotFound
	RatingNotRanked
	RatingOwner
	RatingReady
	RatingInvalid
	RatingAlreadySubmitted
	RatingSubmitted
)

type RatingResult struct {
	Type    RatingResultType
	Average float64
}

func (r *RatingResult) CommonResponse() string {
	switch r.Type {
	case RatingAuthenticationFailed:
		return "auth fail"
	case RatingBeatmapNotFound:
		return "no exist"
	case RatingNotRanked:
		return "not ranked"
	case RatingOwner:
		return "owner"
	case RatingReady:
		return "ok"
	case RatingInvalid:
		return "no"
	case RatingAlreadySubmitted:
		return fmt.Sprintf("alreadyvoted\n%.2f", r.Average)
	case RatingSubmitted:
		return fmt.Sprintf("ok\n%.2f", r.Average)
	default:
		return ""
	}
}

func NewRatingResultWithAverage(ctx *server.Context, result RatingResultType, checksum string) (RatingResult, error) {
	average, err := ctx.State.Ratings.Average(checksum)
	if err != nil {
		return RatingResult{}, fmt.Errorf("calculate beatmap rating average: %w", err)
	}
	return RatingResult{Type: result, Average: average}, nil
}

// /web/osu-rate.php -> Submit a beatmap rating
func OsuRate(ctx *server.Context) {
	result, err := processRating(ctx)
	if err != nil {
		ctx.Logger.Error("Failed to process beatmap rating", "error", err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx.RenderText(http.StatusOK, result.CommonResponse())
}

// /rating/ingame-rate.php -> Submit a beatmap rating (oldest response format)
func IngameRate(ctx *server.Context) {
	result, err := processRating(ctx)
	if err != nil {
		ctx.Logger.Error("Failed to process beatmap rating", "error", err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := result.CommonResponse()
	switch result.Type {
	case RatingAlreadySubmitted:
		response = "alreadyvoted"
	case RatingSubmitted:
		response = fmt.Sprintf("%.2f", result.Average)
	}
	ctx.RenderText(http.StatusOK, response)
}

// /rating/ingame-rate2.php -> Submit a beatmap rating with the average included on subsequent requests
func IngameRate2(ctx *server.Context) {
	result, err := processRating(ctx)
	if err != nil {
		ctx.Logger.Error("Failed to process beatmap rating", "error", err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := result.CommonResponse()
	switch result.Type {
	case RatingAlreadySubmitted:
		response = fmt.Sprintf("alreadyvoted\n%.2f", result.Average)
	case RatingSubmitted:
		response = fmt.Sprintf("%.2f", result.Average)
	}
	ctx.RenderText(http.StatusOK, response)
}

func processRating(ctx *server.Context) (RatingResult, error) {
	user, err := ctx.AuthenticateUserFromQuery("u", "p", true)
	switch {
	case errors.Is(err, server.ErrUserNotFound),
		errors.Is(err, server.ErrInvalidPassword),
		errors.Is(err, server.ErrBanchoPresenceNotFound):
		return RatingResult{Type: RatingAuthenticationFailed}, nil
	case err != nil:
		return RatingResult{}, fmt.Errorf("authenticate user: %w", err)
	}
	user.LatestActivity = time.Now()
	ctx.State.Users.Update(user, "latest_activity")

	checksum := ctx.QueryValue("c")
	beatmap, err := ctx.State.Beatmaps.ByChecksum(checksum, "Beatmapset")
	if err != nil {
		return RatingResult{}, fmt.Errorf("fetch beatmap: %w", err)
	}
	if beatmap == nil {
		return RatingResult{Type: RatingBeatmapNotFound}, nil
	}
	if beatmap.Status <= constants.BeatmapStatusPending {
		return RatingResult{Type: RatingNotRanked}, nil
	}
	if beatmap.Beatmapset.CreatorId != nil && *beatmap.Beatmapset.CreatorId == user.Id {
		return RatingResult{Type: RatingOwner}, nil
	}

	alreadyVoted, err := ctx.State.Ratings.Exists(user.Id, beatmap.Checksum)
	if err != nil {
		return RatingResult{}, fmt.Errorf("check existing beatmap rating: %w", err)
	}
	if alreadyVoted {
		return NewRatingResultWithAverage(ctx, RatingAlreadySubmitted, beatmap.Checksum)
	}

	value, err := ctx.QueryValueIntOptional("v")
	if err != nil {
		return RatingResult{Type: RatingInvalid}, nil
	}
	if value == nil {
		return RatingResult{Type: RatingReady}, nil
	}
	if *value < 0 || *value > 10 {
		return RatingResult{Type: RatingInvalid}, nil
	}

	created, err := ctx.State.Ratings.CreateIfAbsent(&schemas.BeatmapRating{
		UserId:      user.Id,
		SetId:       beatmap.SetId,
		MapChecksum: beatmap.Checksum,
		Rating:      *value,
	})
	if err != nil {
		return RatingResult{}, fmt.Errorf("create beatmap rating: %w", err)
	}
	if !created {
		return NewRatingResultWithAverage(ctx, RatingAlreadySubmitted, beatmap.Checksum)
	}

	err = activity.Submit(
		ctx.State,
		user.Id,
		new(beatmap.Mode),
		constants.ActivityBeatmapRated,
		map[string]any{
			"username":     user.Name,
			"beatmap_id":   beatmap.Id,
			"beatmap_name": beatmap.Name(),
			"rating":       *value,
		},
		false, // should not be sent to #announce
		true,  // should be hidden in user profile
	)
	if err != nil {
		ctx.Logger.Warn(
			"Failed to broadcast beatmap rating activity",
			"user_id", user.Id, "beatmap_id", beatmap.Id, "error", err,
		)
	}

	ctx.Logger.Info(
		"Submitted beatmap rating",
		"user_id", user.Id, "beatmap_id", beatmap.Id, "rating", *value,
	)
	return NewRatingResultWithAverage(ctx, RatingSubmitted, beatmap.Checksum)
}
