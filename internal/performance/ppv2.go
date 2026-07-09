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
	StarRating float64
	// TODO: Add more stuff
}

// TODO: Find pp systems that work with golang & implement them
