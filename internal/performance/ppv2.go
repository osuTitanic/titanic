package performance

import (
	"math"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
)

// IPPv2Service defines an interface for calculating performance points (v2) for scores.
type IPPv2Service interface {
	Available() bool
	CalculatePerformance(score *schemas.Score) (float64, error)
	CalculateDifficulty(beatmapId int, mode constants.Mode, mods constants.Mods) (*DifficultyAttributes, error)
	SetCacheLayer(cache PPv2CacheLayer)
}

type DifficultyAttributes struct {
	Mode       constants.Mode
	IsConvert  bool
	StarRating float64
	Aim        float64
	Speed      float64
	Flashlight float64
	Stamina    float64
	Rhythm     float64
	Color      float64
	Reading    float64
	MaxCombo   uint32
}

// CalculateWeightedPPv2 calculates a user's ppv2 total from scores ordered by descending pp.
func CalculateWeightedPPv2(scores []*schemas.Score) float64 {
	weightedPP := 0.0
	scoreCount := 0

	for _, score := range scores {
		if score == nil {
			continue
		}

		weightedPP += score.PP * math.Pow(0.95, float64(scoreCount))
		scoreCount++
	}
	if scoreCount == 0 {
		return 0
	}

	bonusPP := 416.6667 * (1 - math.Pow(0.9994, float64(scoreCount)))
	return weightedPP + bonusPP
}

// CalculateWeightedAccuracy calculates a user's weighted accuracy from scores ordered by descending pp.
func CalculateWeightedAccuracy(scores []*schemas.Score) float64 {
	weightedAccuracy := 0.0
	scoreCount := 0

	for _, score := range scores {
		if score == nil {
			continue
		}

		weightedAccuracy += score.Acc * math.Pow(0.95, float64(scoreCount))
		scoreCount++
	}
	if scoreCount == 0 {
		return 0
	}

	weightTotal := 20 * (1 - math.Pow(0.95, float64(scoreCount)))
	return weightedAccuracy / weightTotal
}
