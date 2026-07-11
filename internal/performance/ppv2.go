package performance

import (
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
)

// IPPv2Service defines an interface for calculating performance points (v2) for scores.
type IPPv2Service interface {
	CalculatePerformance(score *schemas.Score) (float64, error)
	CalculateDifficulty(beatmapId int, mode constants.Mode, mods constants.Mods) (*DifficultyAttributes, error)
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
