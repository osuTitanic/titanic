//go:build native

package performance

import (
	"math"
	"testing"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
)

func TestPPv2ServiceNative(t *testing.T) {
	provider := &testProvider{}
	service := NewPPv2Service(provider)

	difficulty, err := service.CalculateDifficulty(75, constants.ModeOsu, constants.NoMod)
	if err != nil {
		t.Fatalf("calculate difficulty: %v", err)
	}
	if !provider.closed {
		t.Fatal("beatmap stream was not closed")
	}
	if difficulty.StarRating <= 0 || math.IsNaN(difficulty.StarRating) {
		t.Fatalf("unexpected star rating: %f", difficulty.StarRating)
	}
	provider.closed = false

	pp, err := service.CalculatePerformance(&schemas.Score{
		BeatmapId:  75,
		Mode:       constants.ModeOsu,
		MaxCombo:   314,
		Count300:   194,
		TotalScore: 1_491_676,
	})
	if err != nil {
		t.Fatalf("calculate performance: %v", err)
	}
	if !provider.closed {
		t.Fatal("beatmap stream was not closed")
	}
	if pp <= 0 || math.IsNaN(pp) {
		t.Fatalf("unexpected performance: %f", pp)
	}
}

func TestNativeModAcronyms(t *testing.T) {
	for _, entry := range nativeModAcronyms {
		t.Run(entry.acronym, func(t *testing.T) {
			mods, err := newNativeMods(entry.mod)
			if err != nil {
				t.Fatal(err)
			}
			mods.Close()
		})
	}
}
