package scoring

import (
	"errors"
	"fmt"
	"math"

	"github.com/osuTitanic/titanic/internal/activity"
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
)

func (processor *Processor) broadcastRankHighlight() error {
	if !processor.canBroadcastHighlights() {
		return nil
	}
	submission := processor.submission

	previousRank := submission.PreviousStats.Rank
	currentRank := submission.CurrentStats.Rank
	activityType, isAnnouncement, ok := resolveRankHighlight(previousRank, currentRank)
	if !ok {
		return nil
	}

	data := map[string]any{
		"username":     submission.User.Name,
		"ranks_gained": previousRank - currentRank,
		"rank":         currentRank,
		"mode":         submission.Mode.String(),
	}
	if activityType == constants.ActivityNumberOne {
		data = map[string]any{
			"username": submission.User.Name,
			"mode":     submission.Mode.String(),
		}
	}

	activityErr := activity.Submit(
		processor.context.State,
		submission.UserId,
		&submission.Mode,
		activityType,
		data,
		isAnnouncement,
		false,
	)
	if activityType != constants.ActivityNumberOne {
		return activityErr
	}

	// Player is now #1, send them a small notification :)
	notificationErr := processor.repositories.Notifications.Create(&schemas.Notification{
		UserId: submission.UserId,
		Type:   constants.NotificationTypeAchievement,
		Header: "Welcome to the top!",
		Content: fmt.Sprintf(
			"Congratulations for reaching the #1 global rank in %s. "+
				"Your incredible skill and dedication have set you apart as the absolute best in the game. "+
				"Best of luck on your continued journey at the top!",
			submission.Mode.String(),
		),
	})
	return errors.Join(activityErr, notificationErr)
}

func (processor *Processor) broadcastPerformanceHighlight() error {
	if !processor.canBroadcastHighlights() {
		return nil
	}
	submission := processor.submission

	if submission.StatusPP != constants.ScoreStatusBest {
		// Score is not visible on global rankings
		return nil
	}

	// Get current pp record for mode
	record, err := processor.repositories.Scores.FetchPPRecord(submission.Mode)
	if err != nil {
		return fmt.Errorf("fetch pp record: %w", err)
	}
	if record == nil {
		// No score has ever been set before
		return nil
	}

	// Get player's current top play
	topPlays, err := processor.repositories.Scores.FetchBestRange(
		submission.UserId,
		submission.Mode,
		1, 0,
	)
	if err != nil {
		return fmt.Errorf("fetch player top play: %w", err)
	}

	var topPlay *schemas.Score
	if len(topPlays) > 0 {
		topPlay = topPlays[0]
	}

	// This will return either a pp record to announce or a new
	// top play displayed on the user's profile activity section
	activityType, isAnnouncement, ok := resolvePerformanceHighlight(
		submission.Id,
		record,
		topPlay,
	)
	if !ok {
		return nil
	}

	return activity.Submit(
		processor.context.State,
		submission.UserId,
		&submission.Mode,
		activityType,
		map[string]any{
			"username":   submission.User.Name,
			"beatmap":    submission.Beatmap.Name(),
			"beatmap_id": submission.BeatmapId,
			"pp":         int(math.Round(submission.PP)),
			"mode":       submission.Mode.String(),
		},
		isAnnouncement, // sent to #announce if true
		false,          // should be displayed on their user profile
	)
}

func (processor *Processor) broadcastBeatmapHighlight() error {
	if !processor.canBroadcastHighlights() {
		return nil
	}
	submission := processor.submission

	// Get short-form mods string (i.e. HDHR)
	mods := ""
	if submission.Mods != constants.NoMod {
		mods = submission.Mods.String()
	}

	activityType, isAnnouncement, isHidden := resolveBeatmapHighlight(
		submission.Score,
		submission.NewBeatmapRank,
	)
	activityErr := activity.Submit(
		processor.context.State,
		submission.UserId,
		&submission.Mode,
		activityType,
		map[string]any{
			"username":     submission.User.Name,
			"beatmap":      submission.Beatmap.Name(),
			"beatmap_id":   submission.BeatmapId,
			"beatmap_rank": submission.NewBeatmapRank,
			"mode":         submission.Mode.String(),
			"mods":         mods,
			"pp":           int(math.Round(submission.PP)),
		},
		isAnnouncement,
		isHidden,
	)

	if submission.StatusScore != constants.ScoreStatusBest {
		// Score is not visible on global rankings
		return activityErr
	}
	if submission.NewBeatmapRank != 1 {
		// Score is not #1 on the beatmap
		return activityErr
	}
	if submission.OldBeatmapRank == submission.NewBeatmapRank {
		// User already had #1 on the beatmap
		return activityErr
	}

	topScores, err := processor.repositories.Scores.FetchRangeScores(
		submission.BeatmapId,
		submission.Mode,
		2, 0,
		"User",
	)
	if err != nil {
		return errors.Join(activityErr, fmt.Errorf("fetch displaced beatmap leader: %w", err))
	}
	if len(topScores) <= 1 {
		// No other scores on the beatmap
		return activityErr
	}
	secondPlace := topScores[1]

	if secondPlace.UserId == submission.UserId {
		// The user already had #1 on the beatmap
		// so they never actually "lost" their #1 spot
		return activityErr
	}

	// Show on user profile that the player has
	// lost their #1 spot on the beatmap
	displacedActivityErr := activity.Submit(
		processor.context.State,
		secondPlace.UserId,
		&submission.Mode,
		constants.ActivityLostFirstPlace,
		map[string]any{
			"username":   secondPlace.User.Name,
			"beatmap":    submission.Beatmap.Name(),
			"beatmap_id": submission.BeatmapId,
			"mode":       submission.Mode.String(),
		},
		false,
		false,
	)

	// Fetch the stats of the player who lost their #1 spot
	secondPlaceStats, err := processor.repositories.Stats.ByMode(
		secondPlace.UserId,
		submission.Mode.Value(),
	)
	if err != nil {
		return errors.Join(activityErr, displacedActivityErr, fmt.Errorf("fetch displaced beatmap leader stats: %w", err))
	}
	if secondPlaceStats == nil {
		return errors.Join(activityErr, displacedActivityErr, errors.New("displaced beatmap leader stats missing"))
	}

	// The displaced player has one less #1 score now, update rankings accordingly
	rankingErr := processor.context.State.Rankings.UpdateLeaderScores(
		secondPlaceStats,
		secondPlace.User.Country,
		processor.repositories.Scores,
	)
	return errors.Join(activityErr, displacedActivityErr, rankingErr)
}

func resolveRankHighlight(previousRank, currentRank int) (activityType constants.UserActivity, isAnnouncement bool, ok bool) {
	userUnranked := previousRank <= 0 || currentRank <= 0
	userDeranked := currentRank >= previousRank
	if userUnranked || userDeranked {
		return 0, false, false
	}

	switch {
	case currentRank == 1:
		// Player is now #1
		return constants.ActivityNumberOne, true, true
	case currentRank <= 10:
		// Player has risen to the top 10 or above
		return constants.ActivityRanksGained, true, true
	case previousRank > 100 && currentRank <= 100:
		// Player has risen to the top 100
		return constants.ActivityRanksGained, true, true
	case previousRank > 1000 && currentRank <= 1000:
		// Player has risen to the top 1000
		return constants.ActivityRanksGained, false, true
	default:
		return 0, false, false
	}
}

func resolvePerformanceHighlight(scoreId int64, record, topPlay *schemas.Score) (activityType constants.UserActivity, isAnnouncement bool, ok bool) {
	if record == nil {
		return 0, false, false
	}
	if scoreId == record.Id {
		// Player has set the new pp record
		return constants.ActivityPPRecord, true, true
	}
	if topPlay != nil && scoreId == topPlay.Id {
		// Player got a new top play
		return constants.ActivityTopPlay, false, true
	}
	return 0, false, false
}

func resolveBeatmapHighlight(score *schemas.Score, rank int) (activityType constants.UserActivity, isAnnouncement bool, isHidden bool) {
	visible := score.StatusScore == constants.ScoreStatusBest && rank > 0 && rank <= 1000
	if !visible {
		return constants.ActivityScoreSubmitted, false, true
	}
	return constants.ActivityBeatmapLeaderboardRank, rank <= 4, false
}

func (processor *Processor) canBroadcastHighlights() bool {
	submission := processor.submission
	if submission == nil || submission.Beatmap == nil || submission.User == nil {
		// if this ever happens i'm going to question reality itself
		return false
	}
	return submission.Passed && submission.Id > 0 && !submission.Hidden
}
