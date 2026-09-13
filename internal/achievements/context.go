package achievements

import (
	"sync"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
)

type AchievementContext struct {
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

// Evaluate checks all achievement conditions concurrently, then returns all matches.
func (context *AchievementContext) Evaluate(unlockedFilenames map[string]struct{}) []AchievementDefinition {
	matched := make([]bool, len(Definitions))

	var wg sync.WaitGroup
	for index := range Definitions {
		if _, unlocked := unlockedFilenames[Definitions[index].Filename]; unlocked {
			// Achievement was already unlocked
			continue
		}

		wg.Go(func() {
			matched[index] = Definitions[index].Check(context)
		})
	}
	wg.Wait()

	result := make([]AchievementDefinition, 0)
	for index, unlocked := range matched {
		if !unlocked {
			continue
		}
		result = append(result, Definitions[index])
	}
	return result
}

// TODO: Add context loader / creator
