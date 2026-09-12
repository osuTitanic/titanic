package scoring

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/rankings"
)

func (processor *Processor) finalize() {
	run := func(name string, action func() error) {
		start := time.Now()
		err := action()
		if err != nil {
			processor.AddWarning("%s: %w", name, err)
		}
		elapsed := time.Since(start)

		processor.context.Logger.Debug(
			"Score submission task completed",
			"step", name, "took", elapsed.String(),
		)
	}

	run("synchronize rankings", processor.synchronizeRankings)
	run("synchronize leader rankings", processor.synchronizeLeaderRankings)
	run("update global rank", processor.updateGlobalRank)
	run("update rank history", processor.updateRankHistory)
	run("resolve beatmap rank", processor.resolveNewBeatmapRank)

	// TODO: Check & unlock achievements

	run("populate chart values", processor.setChartValuesAfter)
	run("resolve next overall rank", processor.resolveNextRankOverall)
	run("populate beatmap chart", processor.setBeatmapChartAfter)
	run("upload replay", processor.uploadReplay)
	run("notify bancho", processor.banchoUserUpdateSignal)

	// TODO: Broadcast rank, beatmap, and performance highlights / activity
}

func (processor *Processor) synchronizeRankings() error {
	return processor.context.State.Rankings.Update(
		processor.submission.CurrentStats,
		processor.submission.User.Country,
	)
}

func (processor *Processor) synchronizeLeaderRankings() error {
	return processor.context.State.Rankings.UpdateLeaderScores(
		processor.submission.CurrentStats,
		processor.submission.User.Country,
		processor.repositories.Scores,
	)
}

func (processor *Processor) updateGlobalRank() error {
	stats := processor.submission.CurrentStats

	// Fetch latest rank from redis & apply it to database
	rank, err := processor.context.State.Rankings.GlobalRank(stats.UserId, stats.Mode)
	if err != nil {
		return err
	}

	stats.Rank = rank
	_, err = processor.repositories.Stats.Update(stats, "rank")
	return err
}

func (processor *Processor) updateRankHistory() error {
	if processor.context.State.Config.FrozenRankUpdates {
		return nil
	}
	_, err := processor.repositories.Histories.UpdateRank(
		processor.submission.CurrentStats,
		processor.submission.User.Country,
		processor.context.State.Rankings,
	)
	return err
}

func (processor *Processor) resolveNewBeatmapRank() error {
	if processor.submission.Id <= 0 {
		return nil
	}
	rank, err := processor.repositories.Scores.FetchScoreIndex(processor.submission.Score)
	if err != nil {
		return err
	}
	processor.submission.NewBeatmapRank = rank
	return nil
}

func (processor *Processor) setChartValuesAfter() error {
	stats := processor.submission.CurrentStats
	charts := processor.submission.Charts
	user := processor.submission.User

	accuracy := stats.Accuracy()
	charts.Overall.RankedScore.After = new(stats.Rscore)
	charts.Overall.TotalScore.After = new(stats.Tscore)
	charts.Overall.PlayCount.After = new(stats.Playcount)
	charts.Overall.MaxCombo.After = new(stats.MaxCombo)
	charts.Overall.Accuracy.After = new(accuracy)
	charts.Overall.Rank.After = new(stats.Rank)
	charts.Overall.PP.After = new(stats.PP)

	charts.Overall.BeatmapRanking.Before = new(processor.submission.OldBeatmapRank)
	charts.Overall.BeatmapRanking.After = new(processor.submission.NewBeatmapRank)
	charts.Overall.OnlineScoreId = processor.submission.Id

	// On Titanic, we allow users to set their preferred ranking type on bancho, e.g. PPv1 or Score
	// For score submission, we need to respond with the same rank type that they set on bancho
	preferredRank, preferredValue, err := processor.resolvePreferredRanking()
	if err != nil {
		return err
	}

	// Use the preferred rank
	charts.Overall.Rank.After = new(preferredRank)

	if user.PreferredRanking == constants.RankingTypePPv1 {
		// Use ppv1 over ppv2 if the rank type is set accordingly
		charts.Overall.PP.After = new(preferredValue)
	}
	return nil
}

func (processor *Processor) resolveNextRankOverall() error {
	user := processor.submission.User
	rankType := user.PreferredRanking.Alias()

	difference, aboveUserId, err := processor.context.State.Rankings.PlayerAbove(
		user.Id,
		processor.submission.Mode,
		rankType,
	)
	if errors.Is(err, rankings.ErrNoPlayerAbove) {
		return nil
	}
	if err != nil {
		return err
	}

	aboveUser, err := processor.repositories.Users.ById(aboveUserId)
	if err != nil {
		return err
	}
	if aboveUser == nil {
		return nil
	}
	processor.submission.Charts.Overall.ToNextRank = difference
	processor.submission.Charts.Overall.ToNextRankUser = aboveUser.Name
	return nil
}

func (processor *Processor) setBeatmapChartAfter() error {
	if !processor.submission.Endpoint.IncludesBeatmapRankingChart() {
		return nil
	}

	// Don't display pp awards for loved / qualified maps
	awardsPP := processor.submission.Beatmap.AwardsPP()

	chart := &processor.submission.Charts.Ranking
	oldScore := processor.submission.PersonalBestScore

	if oldScore != nil {
		accuracy := oldScore.Acc * 100
		chart.Rank.Before = new(processor.submission.OldBeatmapRank)
		chart.RankedScore.Before = new(oldScore.TotalScore)
		chart.TotalScore.Before = new(oldScore.TotalScore)
		chart.MaxCombo.Before = new(oldScore.MaxCombo)
		chart.Accuracy.Before = new(accuracy)

		chart.PP.Before = new(oldScore.PP)
		if !awardsPP {
			chart.PP.Before = nil
		}
	}

	accuracy := processor.submission.Acc * 100
	chart.Rank.After = new(processor.submission.NewBeatmapRank)
	chart.RankedScore.After = new(processor.submission.TotalScore)
	chart.TotalScore.After = new(processor.submission.TotalScore)
	chart.MaxCombo.After = new(processor.submission.MaxCombo)
	chart.Accuracy.After = new(accuracy)

	chart.PP.After = new(processor.submission.PP)
	if !awardsPP {
		chart.PP.After = nil
	}

	scoreAbove, err := processor.repositories.Scores.FetchScoreAbove(
		processor.submission.BeatmapId,
		processor.submission.Mode,
		processor.submission.TotalScore,
		"User",
	)
	if err != nil {
		return err
	}
	if scoreAbove == nil || scoreAbove.User == nil {
		return nil
	}
	chart.ToNextRank = scoreAbove.TotalScore - processor.submission.TotalScore
	chart.ToNextRankUser = scoreAbove.User.Name
	return nil
}

func (processor *Processor) uploadReplay() error {
	submission := processor.submission
	if !submission.Passed || len(submission.Replay) == 0 || submission.Id <= 0 {
		return nil
	}
	if submission.StatusPP <= constants.ScoreStatusExited {
		return nil
	}
	if len(submission.Replay) > MaxReplaySize {
		return fmt.Errorf("replay exceeds 10 MB")
	}

	// TODO: Use NewBeatmapRank to determine if replay should be uploaded
	// TODO: Defer upload to go routine
	// TODO: Cache replay

	return processor.context.State.Storage.Save(
		strconv.FormatInt(submission.Id, 10),
		"replays",
		submission.Replay,
	)
}

// TODO: Add a shared helper for sending bancho events
//		 internal/activity could be a good place maybe?

func (processor *Processor) banchoUserUpdateSignal() error {
	type userUpdate struct {
		UserId int            `json:"user_id"`
		Mode   constants.Mode `json:"mode"`
	}
	type event struct {
		Event  string     `json:"event"`
		Args   []any      `json:"args"`
		Kwargs userUpdate `json:"kwargs"`
	}

	payload, err := json.Marshal(event{
		Event: "user_update",
		Args:  []any{},
		Kwargs: userUpdate{
			UserId: processor.submission.UserId,
			Mode:   processor.submission.Mode,
		},
	})
	if err != nil {
		return err
	}

	ctx := processor.context.Request.Context()
	return processor.context.State.Redis.Publish(ctx, "bancho:events", payload).Err()
}
