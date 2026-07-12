//go:build native

package performance

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	osunative "github.com/7mochi/osu-native-go"
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/resources"
	"github.com/osuTitanic/titanic/internal/schemas"
)

type PPv2ServiceNative struct {
	provider resources.BeatmapResourceProvider
	cache    PPv2CacheLayer
}

func NewPPv2Service(provider resources.BeatmapResourceProvider) (IPPv2Service, error) {
	return &PPv2ServiceNative{provider: provider}, nil
}

func (service *PPv2ServiceNative) SetCacheLayer(cache PPv2CacheLayer) {
	service.cache = cache
}

func (service *PPv2ServiceNative) CalculatePerformance(score *schemas.Score) (float64, error) {
	if score == nil {
		return 0, nil
	}

	beatmap, _, err := service.LoadBeatmap(score.BeatmapId)
	if err != nil {
		return 0, err
	}
	defer beatmap.Close()

	ruleset, err := osunative.NewRulesetFromID(score.Mode.Value())
	if err != nil {
		return 0, fmt.Errorf("create ruleset for beatmap %d: %w", score.BeatmapId, err)
	}
	defer ruleset.Close()

	mods, err := newNativeMods(score.Mods)
	if err != nil {
		return 0, fmt.Errorf("create mods for beatmap %d: %w", score.BeatmapId, err)
	}
	defer mods.Close()

	difficultyCalculator, err := osunative.CreateDifficultyCalculator(ruleset, beatmap)
	if err != nil {
		return 0, fmt.Errorf("create difficulty calculator for beatmap %d: %w", score.BeatmapId, err)
	}
	defer difficultyCalculator.Close()

	difficulty, err := difficultyCalculator.Calculate(mods)
	if err != nil {
		return 0, fmt.Errorf("calculate difficulty for beatmap %d: %w", score.BeatmapId, err)
	}

	performanceCalculator, err := osunative.CreatePerformanceCalculator(ruleset)
	if err != nil {
		return 0, fmt.Errorf("create performance calculator for beatmap %d: %w", score.BeatmapId, err)
	}
	defer performanceCalculator.Close()

	attributes, err := performanceCalculator.Calculate(
		ruleset,
		beatmap,
		mods,
		newNativeScoreInfo(score),
		difficulty,
	)
	if err != nil {
		return 0, fmt.Errorf("calculate performance for beatmap %d: %w", score.BeatmapId, err)
	}

	return nativePerformanceTotal(attributes)
}

func (service *PPv2ServiceNative) CalculateDifficulty(beatmapId int, mode constants.Mode, mods constants.Mods) (*DifficultyAttributes, error) {
	beatmap, beatmapMode, err := service.LoadBeatmap(beatmapId)
	if err != nil {
		return nil, err
	}
	defer beatmap.Close()

	ruleset, err := osunative.NewRulesetFromID(mode.Value())
	if err != nil {
		return nil, fmt.Errorf("create ruleset for beatmap %d: %w", beatmapId, err)
	}
	defer ruleset.Close()

	nativeMods, err := newNativeMods(mods)
	if err != nil {
		return nil, fmt.Errorf("create mods for beatmap %d: %w", beatmapId, err)
	}
	defer nativeMods.Close()

	calculator, err := osunative.CreateDifficultyCalculator(ruleset, beatmap)
	if err != nil {
		return nil, fmt.Errorf("create difficulty calculator for beatmap %d: %w", beatmapId, err)
	}
	defer calculator.Close()

	attributes, err := calculator.Calculate(nativeMods)
	if err != nil {
		return nil, fmt.Errorf("calculate difficulty for beatmap %d: %w", beatmapId, err)
	}

	result, err := newDifficultyAttributes(mode, mode != beatmapMode, attributes)
	if err != nil {
		return nil, fmt.Errorf("calculate difficulty for beatmap %d: %w", beatmapId, err)
	}
	return result, nil
}

func (service *PPv2ServiceNative) LoadBeatmap(beatmapId int) (*osunative.Beatmap, constants.Mode, error) {
	if service == nil || service.provider == nil {
		return nil, constants.ModeOsu, fmt.Errorf("load beatmap %d: beatmap resource provider is nil", beatmapId)
	}

	stream, err := service.provider.Osu(beatmapId)
	if err != nil {
		return nil, constants.ModeOsu, fmt.Errorf("load beatmap %d: %w", beatmapId, err)
	}
	if stream == nil {
		return nil, constants.ModeOsu, fmt.Errorf("load beatmap %d: resource provider returned a nil stream", beatmapId)
	}
	defer stream.Close()

	data, err := io.ReadAll(stream)
	if err != nil {
		return nil, constants.ModeOsu, fmt.Errorf("read beatmap %d: %w", beatmapId, err)
	}

	// TODO: Ensure this is decoded as utf-8-bom
	beatmap, err := osunative.NewBeatmapFromText(string(data))
	if err != nil {
		return nil, constants.ModeOsu, fmt.Errorf("parse beatmap %d: %w", beatmapId, err)
	}

	return beatmap, nativeBeatmapMode(data), nil
}

func newNativeMods(mods constants.Mods) (*osunative.ModsCollection, error) {
	collection, err := osunative.NewModsCollection()
	if err != nil {
		return nil, err
	}
	acronyms := []string{"CL"}

	for _, entry := range nativeModAcronyms {
		if !mods.Has(entry.mod) {
			continue
		}
		if entry.mod == constants.DoubleTime && mods.Has(constants.Nightcore) {
			// DTNC moment
			continue
		}
		if entry.mod == constants.SuddenDeath && mods.Has(constants.Perfect) {
			// SDPF moment
			continue
		}
		acronyms = append(acronyms, entry.acronym)
	}

	for _, acronym := range acronyms {
		mod, err := osunative.NewMod(acronym)
		if err != nil {
			collection.Close()
			return nil, fmt.Errorf("create mod %s: %w", acronym, err)
		}
		if err := collection.Add(mod); err != nil {
			mod.Close()
			collection.Close()
			return nil, fmt.Errorf("add mod %s: %w", acronym, err)
		}
	}

	return collection, nil
}

var nativeModAcronyms = []struct {
	mod     constants.Mods
	acronym string
}{
	{constants.NoFail, "NF"},
	{constants.Easy, "EZ"},
	{constants.Hidden, "HD"},
	{constants.HardRock, "HR"},
	{constants.SuddenDeath, "SD"},
	{constants.DoubleTime, "DT"},
	{constants.Relax, "RX"},
	{constants.HalfTime, "HT"},
	{constants.Nightcore, "NC"},
	{constants.Flashlight, "FL"},
	{constants.Autoplay, "AT"},
	{constants.SpunOut, "SO"},
	{constants.Autopilot, "AP"},
	{constants.Perfect, "PF"},
	{constants.Key4, "4K"},
	{constants.Key5, "5K"},
	{constants.Key6, "6K"},
	{constants.Key7, "7K"},
	{constants.Key8, "8K"},
	{constants.FadeIn, "FI"},
	{constants.Random, "RD"},
	{constants.Cinema, "CN"},
	{constants.Target, "TP"},
	{constants.Key9, "9K"},
	{constants.KeyCoop, "CO"},
	{constants.Key1, "1K"},
	{constants.Key3, "3K"},
	{constants.Key2, "2K"},
	{constants.ScoreV2, "SV2"},
	{constants.Mirror, "MR"},
}

func newNativeScoreInfo(score *schemas.Score) *osunative.ScoreInfo {
	totalScore := score.TotalScore

	// TODO: Is this the right mapping????
	return &osunative.ScoreInfo{
		Accuracy:         score.Accuracy(),
		MaxCombo:         score.MaxCombo,
		CountMiss:        score.CountMiss,
		CountMeh:         score.Count50,
		CountOk:          score.Count100,
		CountGood:        score.CountKatu,
		CountGreat:       score.Count300,
		CountPerfect:     score.CountGeki,
		LegacyTotalScore: &totalScore,
	}
}

func newDifficultyAttributes(mode constants.Mode, isConvert bool, attributes any) (*DifficultyAttributes, error) {
	result := &DifficultyAttributes{Mode: mode, IsConvert: isConvert}

	switch value := attributes.(type) {
	case *osunative.OsuDifficultyAttributes:
		result.StarRating = value.StarRating
		result.Aim = value.AimDifficulty
		result.Speed = value.SpeedDifficulty
		result.Flashlight = value.FlashlightDifficulty
		result.MaxCombo = nativeMaxCombo(value.MaxCombo)
	case *osunative.TaikoDifficultyAttributes:
		result.StarRating = value.StarRating
		result.Stamina = value.StaminaDifficulty
		result.Rhythm = value.RhythmDifficulty
		result.Color = value.ColourDifficulty
		result.Reading = value.ReadingDifficulty
		result.MaxCombo = nativeMaxCombo(value.MaxCombo)
	case *osunative.CatchDifficultyAttributes:
		result.StarRating = value.StarRating
		result.MaxCombo = nativeMaxCombo(value.MaxCombo)
	case *osunative.ManiaDifficultyAttributes:
		result.StarRating = value.StarRating
		result.MaxCombo = nativeMaxCombo(value.MaxCombo)
	default:
		return nil, fmt.Errorf("unexpected difficulty attributes %T", attributes)
	}

	return result, nil
}

func nativePerformanceTotal(attributes any) (float64, error) {
	// TODO: Check if we can `.(PerformanceAttributes).Total` here
	switch value := attributes.(type) {
	case *osunative.OsuPerformanceAttributes:
		return value.Total, nil
	case *osunative.TaikoPerformanceAttributes:
		return value.Total, nil
	case *osunative.CatchPerformanceAttributes:
		return value.Total, nil
	case *osunative.ManiaPerformanceAttributes:
		return value.Total, nil
	default:
		return 0, fmt.Errorf("unexpected performance attributes %T", attributes)
	}
}

func nativeMaxCombo(value int32) uint32 {
	if value <= 0 {
		return 0
	}
	return uint32(value)
}

func nativeBeatmapMode(data []byte) constants.Mode {
	// i don't like this, but this is what we have to do i suppose...

	mode := constants.ModeOsu
	scanner := bufio.NewScanner(bytes.NewReader(data))
	inGeneral := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Check for new sections, specifically [General] where the mode lies
		if strings.HasPrefix(line, "[") {
			if inGeneral {
				break
			}
			inGeneral = line == "[General]"
			continue
		}
		if !inGeneral {
			continue
		}

		value, found := strings.CutPrefix(line, "Mode:")
		if !found {
			continue
		}

		// We found our mode
		parsed, err := strconv.Atoi(strings.TrimSpace(value))
		if err == nil {
			mode = constants.Mode(parsed)
		}
		break
	}
	return mode
}
