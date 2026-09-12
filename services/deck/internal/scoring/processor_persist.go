package scoring

import (
	"fmt"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/performance"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
)

func (processor *Processor) persist(transaction *state.Repositories) (ResultType, error) {
	submission := processor.submission
	if submission.User == nil {
		return ResultAccepted, fmt.Errorf("submission user missing")
	}
	if submission.Beatmap == nil {
		return ResultAccepted, fmt.Errorf("submission beatmap missing")
	}
	submission.UserId = submission.User.Id
	submission.BeatmapId = submission.Beatmap.Id

	stats, err := transaction.Stats.ByModeWithLock(submission.UserId, submission.Mode)
	if err != nil {
		return ResultAccepted, fmt.Errorf("lock stats for user %d mode %d: %w", submission.UserId, submission.Mode, err)
	}
	if stats == nil {
		return ResultAccepted, fmt.Errorf("stats not found for user %d mode %d", submission.UserId, submission.Mode)
	}
	previousStats := *stats
	submission.PreviousStats = &previousStats
	submission.CurrentStats = stats

	awardsScore := submission.Beatmap.Status > constants.BeatmapStatusPending

	// Populate "Before" section of charts
	processor.setChartValuesBefore()

	if awardsScore {
		// Recalculate score statuses & save score to database
		if err := processor.persistScore(transaction); err != nil {
			return ResultAccepted, err
		}
	}

	// Update user & beatmap statistics
	if err := processor.updateSubmissionStatistics(transaction, stats, awardsScore); err != nil {
		return ResultAccepted, err
	}

	// TODO: Eventually we want to submit scores on unranked beatmaps (temporarily) as well
	//		 This check would become obsolete with that change
	if !awardsScore {
		return ResultBeatmapUnavailable, nil
	}
	if !processor.context.State.Config.AllowRelax && submission.Relaxing() {
		return ResultRejected, nil
	}
	return ResultAccepted, nil
}

func (processor *Processor) persistScore(transaction *state.Repositories) error {
	submission := processor.submission

	performanceBest, scoreBest, err := transaction.Scores.FetchPersonalBests(
		submission.BeatmapId,
		submission.UserId,
		submission.Mode,
	)
	if err != nil {
		return fmt.Errorf("load personal bests: %w", err)
	}
	submission.PersonalBestPP = performanceBest
	submission.PersonalBestScore = scoreBest

	if scoreBest != nil {
		submission.OldBeatmapRank, err = transaction.Scores.FetchScoreIndexById(
			scoreBest.Id,
			submission.BeatmapId,
			submission.Mode,
		)
		if err != nil {
			return fmt.Errorf("resolve previous beatmap rank: %w", err)
		}
	}

	if err := processor.assignPerformanceStatus(transaction, performanceBest); err != nil {
		return err
	}
	if err := processor.assignScoreStatus(transaction, scoreBest); err != nil {
		return err
	}

	if err := transaction.Scores.Create(submission.Score); err != nil {
		return fmt.Errorf("create score: %w", err)
	}
	return nil
}

func (processor *Processor) updateSubmissionStatistics(transaction *state.Repositories, stats *schemas.Stats, awardsScore bool) error {
	submission := processor.submission

	// Update beatmap statistics
	playcount, passcount, err := transaction.Beatmaps.IncrementPlayCounts(
		submission.BeatmapId,
		submission.Passed,
	)
	if err != nil {
		return fmt.Errorf("update beatmap play counts: %w", err)
	}
	submission.Beatmap.Playcount = playcount
	submission.Beatmap.Passcount = passcount
	submission.Charts.Beatmap.BeatmapPlaycount = playcount
	submission.Charts.Beatmap.BeatmapPasscount = passcount

	// Update user statistics
	stats.Playcount++
	stats.Playtime += submission.ElapsedTime()
	stats.Tscore += submission.TotalScore
	stats.TotalHits += submission.TotalHits()

	if awardsScore && submission.HasPersonalBest() && submission.MaxCombo > stats.MaxCombo {
		stats.MaxCombo = submission.MaxCombo
	}

	// Update "Play History" graph on the user profile
	if err := transaction.Histories.UpdatePlays(submission.UserId, submission.Mode); err != nil {
		return fmt.Errorf("update play history: %w", err)
	}

	// Either increment or create the beatmap plays entry
	if err := transaction.Plays.Increment(&schemas.BeatmapPlays{
		UserId:      submission.UserId,
		BeatmapId:   submission.BeatmapId,
		SetId:       submission.Beatmap.SetId,
		BeatmapFile: submission.Beatmap.Filename,
		Count:       1,
	}); err != nil {
		return fmt.Errorf("update beatmap plays: %w", err)
	}

	bestScores, err := transaction.Scores.FetchBest(
		submission.UserId,
		submission.Mode,
		!processor.context.State.Config.ApprovedMapRewards,
	)
	if err != nil {
		return fmt.Errorf("load performance bests: %w", err)
	}

	stats.PP = performance.CalculateWeightedPPv2(bestScores)
	stats.Acc = performance.CalculateWeightedAccuracy(bestScores)

	if !processor.context.State.Config.FrozenPPv1Updates {
		stats.PPv1 = processor.context.State.PPv1.CalculateWeightFromScores(bestScores)
	}

	scoreBests, err := transaction.Scores.FetchBestByScore(
		submission.UserId,
		submission.Mode,
	)
	if err != nil {
		return fmt.Errorf("load score bests: %w", err)
	}

	// Recalculate ranked score through score pb's
	// TODO: Refactor this into a single database query
	stats.Rscore = 0
	for _, best := range scoreBests {
		stats.Rscore += best.TotalScore
	}

	gradeCounts, err := transaction.Scores.FetchGradeCounts(submission.UserId, submission.Mode)
	if err != nil {
		return fmt.Errorf("load grade counts: %w", err)
	}
	stats.CountXH = gradeCounts[constants.GradeXH]
	stats.CountX = gradeCounts[constants.GradeX]
	stats.CountSH = gradeCounts[constants.GradeSH]
	stats.CountS = gradeCounts[constants.GradeS]
	stats.CountA = gradeCounts[constants.GradeA]
	stats.CountB = gradeCounts[constants.GradeB]
	stats.CountC = gradeCounts[constants.GradeC]
	stats.CountD = gradeCounts[constants.GradeD]

	_, err = transaction.Stats.Update(
		stats,
		"tscore",
		"rscore",
		"pp",
		"ppv1",
		"playcount",
		"playtime",
		"acc",
		"max_combo",
		"total_hits",
		"xh_count",
		"x_count",
		"sh_count",
		"s_count",
		"a_count",
		"b_count",
		"c_count",
		"d_count",
	)
	if err != nil {
		return fmt.Errorf("update stats: %w", err)
	}
	return nil
}

func (processor *Processor) assignPerformanceStatus(transaction *state.Repositories, currentBest *schemas.Score) error {
	score := processor.submission.Score

	if !processor.context.State.Config.AllowRelax && score.Relaxing() {
		score.StatusPP = constants.ScoreStatusHidden
		score.Hidden = true
		return nil
	}

	if score.Relaxing() {
		// On Titanic, we don't award pp for rx/ap scores
		// However, they are still placed on leaderboards
		if processor.submission.Passed {
			score.StatusPP = constants.ScoreStatusSubmitted
		} else {
			score.StatusPP = constants.ScoreStatusExited
		}
		return nil
	}

	if !processor.submission.Passed {
		// Score was a fail or user exited the map
		if processor.submission.Exited {
			score.StatusPP = constants.ScoreStatusExited
		} else {
			score.StatusPP = constants.ScoreStatusFailed
		}
		return nil
	}

	if currentBest == nil {
		// User has never set a score on this map before
		score.StatusPP = constants.ScoreStatusBest
		return nil
	}

	isBetter := score.ComparePerformance(currentBest)

	if isBetter {
		// We have a new personal best!
		// Demote the previous personal best
		currentBest.StatusPP = constants.ScoreStatusSubmitted

		if currentBest.Mods != score.Mods {
			// The previous pb is now a pb with the used mod combination
			// If there's another mods pb, that was worse than this
			// score, we need to demote it as well
			currentBest.StatusPP = constants.ScoreStatusMods

			modsBest, err := transaction.Scores.FetchModsPerformanceBest(
				score.BeatmapId, score.UserId,
				score.Mode, score.Mods,
			)
			if err != nil {
				return fmt.Errorf("fetch mods performance pb: %w", err)
			}

			if modsBest != nil && modsBest.Id != currentBest.Id {
				modsBest.StatusPP = constants.ScoreStatusSubmitted
				if _, err := transaction.Scores.Update(modsBest, "status"); err != nil {
					return fmt.Errorf("demote mods performance pb: %w", err)
				}
			}
		}

		if _, err := transaction.Scores.Update(currentBest, "status"); err != nil {
			return fmt.Errorf("demote previous performance pb: %w", err)
		}
		score.StatusPP = constants.ScoreStatusBest
		return nil
	}

	// Score was not a pb, sadge :(
	// However, maybe it was a pb with the used mod combination? :eyes:

	if currentBest.Mods == score.Mods {
		// The current pb and this score use the same mod combo
		// -> Score is not a mods pb
		score.StatusPP = constants.ScoreStatusSubmitted
		return nil
	}

	modsBest, err := transaction.Scores.FetchModsPerformanceBest(
		score.BeatmapId, score.UserId,
		score.Mode, score.Mods,
	)
	if err != nil {
		return fmt.Errorf("load mod performance best: %w", err)
	}
	if modsBest == nil {
		// There was no other score with this mod combo
		// -> Score is now a mods pb
		score.StatusPP = constants.ScoreStatusMods
		return nil
	}

	isBetterWithMods := score.ComparePerformance(modsBest)
	if !isBetterWithMods {
		// Older pb with mods was better than current score
		// -> Score is not a mods pb
		score.StatusPP = constants.ScoreStatusSubmitted
		return nil
	}

	// Let's gooo, we have a new mods pb
	// -> Demote previous mods pb
	modsBest.StatusPP = constants.ScoreStatusSubmitted
	if _, err := transaction.Scores.Update(modsBest, "status"); err != nil {
		return fmt.Errorf("demote previous mod performance best: %w", err)
	}

	score.StatusPP = constants.ScoreStatusMods
	return nil
}

func (processor *Processor) assignScoreStatus(transaction *state.Repositories, currentBest *schemas.Score) error {
	score := processor.submission.Score

	if !processor.context.State.Config.AllowRelax && score.Relaxing() {
		score.StatusScore = constants.ScoreStatusHidden
		score.Hidden = true
		return nil
	}

	if !processor.submission.Passed {
		// Score was a fail or user exited the map
		if processor.submission.Exited {
			score.StatusScore = constants.ScoreStatusExited
		} else {
			score.StatusScore = constants.ScoreStatusFailed
		}
		return nil
	}

	if currentBest == nil {
		// User has never set a score on this map before
		score.StatusScore = constants.ScoreStatusBest
		return nil
	}

	isBetter := score.CompareScore(currentBest)

	if isBetter {
		// We have a new personal best!
		// Demote the previous personal best
		currentBest.StatusScore = constants.ScoreStatusSubmitted

		if currentBest.Mods != score.Mods {
			// The previous pb is now a pb with the used mod combination
			// If there's another mods pb, that was worse than this
			// score, we need to demote it as well
			currentBest.StatusScore = constants.ScoreStatusMods

			modsBest, err := transaction.Scores.FetchModsScoreBest(
				score.BeatmapId, score.UserId,
				score.Mode, score.Mods,
			)
			if err != nil {
				return fmt.Errorf("load displaced mod score best: %w", err)
			}

			if modsBest != nil && modsBest.Id != currentBest.Id {
				modsBest.StatusScore = constants.ScoreStatusSubmitted
				if _, err := transaction.Scores.Update(modsBest, "status_score"); err != nil {
					return fmt.Errorf("demote displaced mod score best: %w", err)
				}
			}
		}

		if _, err := transaction.Scores.Update(currentBest, "status_score"); err != nil {
			return fmt.Errorf("demote previous score best: %w", err)
		}
		score.StatusScore = constants.ScoreStatusBest
		return nil
	}

	// Score was not a pb, sadge :(
	// However, maybe it was a pb with the used mod combination? :eyes:

	if currentBest.Mods == score.Mods {
		// The current pb and this score use the same mod combo
		// -> Score is not a mods pb
		score.StatusScore = constants.ScoreStatusSubmitted
		return nil
	}

	modsBest, err := transaction.Scores.FetchModsScoreBest(
		score.BeatmapId, score.UserId,
		score.Mode, score.Mods,
	)
	if err != nil {
		return fmt.Errorf("load mod score best: %w", err)
	}
	if modsBest == nil {
		// There was no other score with this mod combo
		// -> Score is now a mods pb
		score.StatusScore = constants.ScoreStatusMods
		return nil
	}

	isBetterWithMods := score.CompareScore(modsBest)
	if !isBetterWithMods {
		// Older pb with mods was better than current score
		// -> Score is not a mods pb
		score.StatusScore = constants.ScoreStatusSubmitted
		return nil
	}

	// Let's gooo, we have a new mods pb
	// -> Demote previous mods pb
	modsBest.StatusScore = constants.ScoreStatusSubmitted
	if _, err := transaction.Scores.Update(modsBest, "status_score"); err != nil {
		return fmt.Errorf("demote previous mod score best: %w", err)
	}
	score.StatusScore = constants.ScoreStatusMods
	return nil
}

func (processor *Processor) setChartValuesBefore() {
	stats := processor.submission.PreviousStats
	charts := processor.submission.Charts
	user := processor.submission.User

	accuracy := stats.Accuracy()
	charts.Overall.RankedScore.Before = new(stats.Rscore)
	charts.Overall.TotalScore.Before = new(stats.Tscore)
	charts.Overall.PlayCount.Before = new(stats.Playcount)
	charts.Overall.MaxCombo.Before = new(stats.MaxCombo)
	charts.Overall.Accuracy.Before = new(accuracy)
	charts.Overall.Rank.Before = new(stats.Rank)
	charts.Overall.PP.Before = new(stats.PP)

	// On Titanic, we allow users to set their preferred ranking type on bancho, e.g. PPv1 or Score
	// For score submission, we need to respond with the same rank type that they set on bancho
	preferredRank, preferredValue, err := processor.resolvePreferredRanking()
	if err != nil {
		processor.AddWarning("resolve previous preferred ranking: %w", err)
		return
	}

	// Use the preferred rank
	charts.Overall.Rank.Before = new(preferredRank)

	if user.PreferredRanking == constants.RankingTypePPv1 {
		// Use ppv1 over ppv2 if the rank type is set accordingly
		charts.Overall.PP.Before = new(preferredValue)
	}
}

func (processor *Processor) resolvePreferredRanking() (rank int, value float64, err error) {
	user := processor.submission.User
	mode := processor.submission.Mode

	rankType := user.PreferredRanking.Alias()
	rankKey := processor.context.State.Rankings.RankingKey(mode, rankType, nil)

	rank, err = processor.context.State.Rankings.Rank(
		user.Id, mode,
		rankType, nil,
	)
	if err != nil {
		return 0, 0, err
	}

	value, err = processor.context.State.Rankings.ScoreByKey(
		rankKey,
		user.Id,
	)
	return rank, value, err
}
