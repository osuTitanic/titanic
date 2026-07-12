//go:build !rosu && !native

package performance

import (
	"errors"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/resources"
	"github.com/osuTitanic/titanic/internal/schemas"
)

type PPv2ServiceDummy struct{}

func NewPPv2Service(provider resources.BeatmapResourceProvider) (IPPv2Service, error) {
	return &PPv2ServiceDummy{}, nil
}

func (service *PPv2ServiceDummy) CalculatePerformance(score *schemas.Score) (float64, error) {
	return 0, errors.New("ppv2 service is not available")
}

func (service *PPv2ServiceDummy) CalculateDifficulty(beatmapId int, mode constants.Mode, mods constants.Mods) (*DifficultyAttributes, error) {
	return nil, errors.New("ppv2 service is not available")
}

func (service *PPv2ServiceDummy) SetCacheLayer(cache PPv2CacheLayer) { /* stub */ }
