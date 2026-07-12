//go:build !rosu && !native

package performance

import (
	"errors"

	"github.com/osuTitanic/titanic/internal/resources"
)

func NewPPv2Service(provider resources.BeatmapResourceProvider) (IPPv2Service, error) {
	return nil, errors.New("a pp system has to be provided through a compile flag")
}
