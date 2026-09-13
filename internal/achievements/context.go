package achievements

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/repositories"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
)

type AchievementContext struct {
	UnlockedFilenames map[string]struct{}

	Score       *schemas.Score
	StatsByMode map[constants.Mode]*schemas.Stats

	// Used by the four profile ranking achievements
	GlobalRank int
	// Used by the Beatmap Pack achievements
	CompletedBeatmapsetIds map[int]struct{}
	// Used by S-Ranker to check the player's five latest scores
	RecentScores []*schemas.Score
	// Used by Quick Draw to identify the first leaderboard score
	BeatmapLeaderboardSize int
	// Used by Obsessed to count retries within the last day
	SameMapPlaysLastDay int
	// Used by Most Improved to find a recent D rank on this beatmap
	HadDOnMapLastDay bool
}

// Evaluate checks all achievement conditions, then returns all matches.
func (context *AchievementContext) Evaluate() (result []AchievementDefinition) {
	for _, definition := range Definitions {
		if _, unlocked := context.UnlockedFilenames[definition.Filename]; unlocked {
			// Achievement was already unlocked
			continue
		}
		if !definition.Check(context) {
			// Achievement was not unlocked
			continue
		}
		result = append(result, definition)
	}
	return result
}

// LoadContext fetches the stuff needed to evaluate achievements for a score
func LoadContext(repos *state.Repositories, score *schemas.Score, currentStats *schemas.Stats) (*AchievementContext, error) {
	context := &AchievementContext{
		Score:             score,
		StatsByMode:       make(map[constants.Mode]*schemas.Stats, len(constants.Modes)),
		UnlockedFilenames: make(map[string]struct{}),
	}
	if currentStats != nil {
		context.StatsByMode[currentStats.Mode] = currentStats
		context.GlobalRank = currentStats.Rank
	}

	err := loadConcurrently(
		// Load already unlocked achievements
		func() error {
			unlocked, err := repos.Achievements.ManyByUserId(score.UserId)
			if err != nil {
				return fmt.Errorf("load unlocked achievements: %w", err)
			}
			for _, achievement := range unlocked {
				if achievement != nil {
					context.UnlockedFilenames[achievement.Filename] = struct{}{} // go doesn't have sets mannnnn
				}
			}
			return nil
		},
		// Load other user stats
		func() error {
			stats, err := repos.Stats.ManyByUserId(score.UserId)
			if err != nil {
				return fmt.Errorf("load user stats: %w", err)
			}
			for _, modeStats := range stats {
				if currentStats == nil || modeStats.Mode != currentStats.Mode {
					context.StatsByMode[modeStats.Mode] = modeStats
				}
			}
			return nil
		},
		// Fetch 5 most recent scores for S-Ranker
		func() error {
			var err error
			context.RecentScores, err = repos.Scores.FetchRecentByUser(
				score.UserId, score.Mode,
				5, constants.ScoreStatusHidden,
			)
			if err != nil {
				return fmt.Errorf("load recent scores: %w", err)
			}
			return nil
		},
		// Fetch set IDs of all completed sets for beatmap packs
		func() error {
			// TODO: Only query beatmap pack IDs
			beatmapsetIds, err := repos.Scores.FetchCompletedBeatmapsetIds(score.UserId)
			if err != nil {
				return fmt.Errorf("load completed beatmap sets: %w", err)
			}

			context.CompletedBeatmapsetIds = make(map[int]struct{}, len(beatmapsetIds))
			for _, beatmapsetId := range beatmapsetIds {
				context.CompletedBeatmapsetIds[beatmapsetId] = struct{}{}
			}
			return nil
		},
		// Fetch data for Most Improved & Obsessed
		func() error {
			var err error
			context.SameMapPlaysLastDay, context.HadDOnMapLastDay, err = repos.Scores.FetchAchievementMapActivity(
				score.UserId,
				score.BeatmapId,
				score.Mode,
				time.Now().UTC().Add(-24*time.Hour),
			)
			if err != nil {
				return fmt.Errorf("load recent beatmap activity: %w", err)
			}
			return nil
		},
		// Fetch leaderboard score count for Quick Draw
		func() error {
			var err error
			context.BeatmapLeaderboardSize, err = repos.Scores.FetchLeaderboardCount(
				repositories.BeatmapLeaderboardFilter{
					BeatmapId: score.BeatmapId,
					Mode:      score.Mode,
				},
			)
			if err != nil {
				return fmt.Errorf("load beatmap leaderboard: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return context, nil
}

func loadConcurrently(loaders ...func() error) error {
	errorsChannel := make(chan error, len(loaders))

	var wg sync.WaitGroup
	for _, load := range loaders {
		wg.Go(func() {
			if err := load(); err != nil {
				errorsChannel <- err
			}
		})
	}
	wg.Wait()
	close(errorsChannel)

	var loadErrors []error
	for err := range errorsChannel {
		loadErrors = append(loadErrors, err)
	}
	return errors.Join(loadErrors...)
}
