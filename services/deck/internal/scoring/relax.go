package scoring

import (
	"fmt"
	"math"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
)

// TODO: We probably want to move this out of internal/scoring, e.g. for recalculations

var totalScoreModMultipliers = [...]struct {
	mod        constants.Mods
	multiplier float64
}{
	{constants.Easy, 0.50},
	{constants.NoFail, 0.50},
	{constants.HalfTime, 0.30},
	{constants.HardRock, 1.06},
	{constants.Hidden, 1.06},
	{constants.DoubleTime, 1.12},
	{constants.Nightcore, 1.12},
	{constants.Flashlight, 1.12},
	{constants.Relax, 0.30},
	{constants.Autopilot, 0.30},
}

// EstimateTotalScore calulates an estimation of the score's total score, used for relax/autopilot
func EstimateTotalScore(score *schemas.Score) (int64, error) {
	if score == nil {
		return 0, fmt.Errorf("estimate total score: score is nil")
	}

	switch score.Mode {
	case constants.ModeOsu:
		return estimateOsuTotalScore(score)
	case constants.ModeCatch:
		return score.TotalScore / 4, nil
	case constants.ModeTaiko, constants.ModeMania:
		return score.TotalScore, nil
	default:
		return score.TotalScore, nil
	}
}

func estimateOsuTotalScore(score *schemas.Score) (int64, error) {
	totalHits := score.Count300 + score.Count100 + score.Count50 + score.CountMiss
	if totalHits <= 0 {
		return 0, nil
	}
	if score.Beatmap == nil {
		return 0, fmt.Errorf("estimate total score: beatmap is nil")
	}
	if score.Beatmap.TotalLength <= 0 {
		return 0, fmt.Errorf("estimate total score: invalid beatmap length %d", score.Beatmap.TotalLength)
	}
	if score.Beatmap.MaxCombo <= 0 {
		return 0, fmt.Errorf("estimate total score: invalid beatmap max combo %d", score.Beatmap.MaxCombo)
	}

	originalScore := float64(score.Count300*300 + score.Count100*100 + score.Count50*50)

	// Account for slider ticks inflating score
	sliderFactor := float64(score.Beatmap.MaxCombo) / float64(totalHits)
	averageHit := (300*float64(score.Count300)/float64(totalHits) +
		100*float64(score.Count100)/float64(totalHits) +
		50*float64(score.Count50)/float64(totalHits)) / sliderFactor

	difficultyMultiplier := estimateDifficultyMultiplier(score.Beatmap, totalHits)
	modMultiplier := estimateModMultiplier(score.Mods.Normalize())

	comboCount := float64(max(score.MaxCombo-1, 0))
	comboSum := comboCount * (comboCount + 1) / 2
	comboScore := averageHit * (comboCount + comboSum*difficultyMultiplier*modMultiplier/25)

	return int64(math.RoundToEven(originalScore + comboScore)), nil
}

func estimateDifficultyMultiplier(beatmap *schemas.Beatmap, totalHits int) float64 {
	lengthFactor := float64(totalHits) / float64(beatmap.TotalLength) * 8
	lengthFactor = min(16, max(0, lengthFactor))
	difficulty := beatmap.HP + beatmap.CS + beatmap.OD + lengthFactor
	return math.RoundToEven(difficulty / 38 * 5)
}

func estimateModMultiplier(mods constants.Mods) float64 {
	multiplier := 1.0
	for _, entry := range totalScoreModMultipliers {
		if mods.Has(entry.mod) {
			multiplier *= entry.multiplier
		}
	}
	return multiplier
}
