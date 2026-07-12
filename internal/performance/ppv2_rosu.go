//go:build rosu

package performance

import (
	"fmt"
	"io"

	rosu "github.com/calemy/rosu-pp-go"
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/resources"
	"github.com/osuTitanic/titanic/internal/schemas"
)

type PPv2ServiceRosu struct {
	provider resources.BeatmapResourceProvider
}

func NewPPv2Service(provider resources.BeatmapResourceProvider) (IPPv2Service, error) {
	return &PPv2ServiceRosu{provider: provider}, nil
}

func (service *PPv2ServiceRosu) CalculatePerformance(score *schemas.Score) (float64, error) {
	if score == nil {
		return 0, nil
	}
	// TODO: rosu-pp-ffi doesn't support converts at the moment

	beatmap, err := service.LoadBeatmap(score.BeatmapId)
	if err != nil {
		return 0, err
	}
	defer beatmap.Free()

	mods := rosu.ModsFromBits(uint32(score.Mods))
	defer mods.Free()

	calculator := beatmap.Calculator()
	calculator.
		Mods(mods).
		Lazer(false).
		Combo(toUint32(score.MaxCombo)).
		Misses(toUint32(score.CountMiss)).
		N300(toUint32(score.Count300)).
		N100(toUint32(score.Count100)).
		N50(toUint32(score.Count50)).
		Geki(toUint32(score.CountGeki)).
		Katu(toUint32(score.CountKatu)).
		LegacyScore(toUint32FromInt64(score.TotalScore))

	attributes, err := calculator.CalculateSafe()
	if err != nil {
		return 0, fmt.Errorf("calculate performance for beatmap %d: %w", score.BeatmapId, err)
	}

	return attributes.PP, nil
}

func (service *PPv2ServiceRosu) CalculateDifficulty(beatmapId int, mode constants.Mode, mods constants.Mods) (*DifficultyAttributes, error) {
	beatmap, err := service.LoadBeatmap(beatmapId)
	if err != nil {
		return nil, err
	}
	defer beatmap.Free()

	if !beatmap.Safe() {
		return nil, fmt.Errorf("calculate difficulty for beatmap %d: beatmap is too suspicious", beatmapId)
	}

	rosuMods := rosu.ModsFromBits(uint32(mods))
	defer rosuMods.Free()

	calculator := rosu.NewDifficulty()
	attributes, err := calculator.
		Mods(rosuMods).
		Lazer(false).
		CalculateSafe(*beatmap)
	if err != nil {
		return nil, fmt.Errorf("calculate difficulty for beatmap %d: %w", beatmapId, err)
	}

	return &DifficultyAttributes{
		Mode:       constants.Mode(attributes.Mode),
		IsConvert:  attributes.IsConvert,
		StarRating: attributes.Stars,
		Aim:        attributes.Aim,
		Speed:      attributes.Speed,
		Flashlight: attributes.Flashlight,
		Stamina:    attributes.Stamina,
		Rhythm:     attributes.Rhythm,
		Color:      attributes.Color,
		Reading:    attributes.Reading,
		MaxCombo:   attributes.MaxCombo,
	}, nil
}

func (service *PPv2ServiceRosu) LoadBeatmap(beatmapId int) (*rosu.Beatmap, error) {
	if service == nil || service.provider == nil {
		return nil, fmt.Errorf("load beatmap %d: beatmap resource provider is nil", beatmapId)
	}

	stream, err := service.provider.Osu(beatmapId)
	if err != nil {
		return nil, fmt.Errorf("load beatmap %d: %w", beatmapId, err)
	}
	if stream == nil {
		return nil, fmt.Errorf("load beatmap %d: resource provider returned a nil stream", beatmapId)
	}
	defer stream.Close()

	data, err := io.ReadAll(stream)
	if err != nil {
		return nil, fmt.Errorf("read beatmap %d: %w", beatmapId, err)
	}

	beatmap, err := rosu.BeatmapFromBytes(data)
	if err != nil {
		return nil, fmt.Errorf("parse beatmap %d: %w", beatmapId, err)
	}

	return beatmap, nil
}

func toUint32(value int) uint32 {
	if value <= 0 {
		return 0
	}
	// ^uint32(0) is the maximum possible uint32 value
	// Clamp values larger than the maximum uint32 value to the max Uint32
	if uint64(value) > uint64(^uint32(0)) {
		return ^uint32(0)
	}
	return uint32(value)
}

func toUint32FromInt64(value int64) uint32 {
	if value <= 0 {
		return 0
	}
	// ^uint32(0) is the maximum possible uint32 value
	// Clamp values larger than the maximum uint32 value to the max Uint32
	if uint64(value) > uint64(^uint32(0)) {
		return ^uint32(0)
	}
	return uint32(value)
}
