package bss

import (
	"crypto/md5"
	"errors"
	"fmt"
	"math"
	"path"
	"strings"

	"github.com/Lekuruu/gosu/pkg/beatmaps"
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/performance"
	"github.com/osuTitanic/titanic/internal/schemas"
)

// AssignBeatmapFile reads a prepared beatmap from the submission filesystem,
// parses it & applies its metadata to the target beatmap.
func (submission *SubmissionContext) AssignBeatmapFile(beatmapId int, filename string) (*PreparedBeatmap, []byte, error) {
	if submission.FS == nil {
		return nil, nil, errors.New("bss: filesystem not set")
	}
	if submission.IsCanceled() {
		return nil, nil, submission.Context.Err()
	}

	if err := validatePackageFilename(filename); err != nil {
		return nil, nil, err
	}
	if !strings.EqualFold(path.Ext(filename), ".osu") {
		return nil, nil, fmt.Errorf("bss: invalid beatmap filename %q", filename)
	}

	contents, err := submission.readFile(filename, MaxBeatmapFileSize)
	if err != nil {
		return nil, nil, fmt.Errorf("bss: read beatmap file %q: %w", filename, err)
	}
	if submission.IsCanceled() {
		return nil, nil, submission.Context.Err()
	}

	beatmap, err := submission.AssignBeatmapContent(beatmapId, filename, contents)
	if err != nil {
		return nil, nil, err
	}
	return beatmap, contents, nil
}

// AssignBeatmapContent parses the given beatmap contents
// & applies it to the matching prepared beatmap.
func (submission *SubmissionContext) AssignBeatmapContent(beatmapId int, filename string, contents []byte) (*PreparedBeatmap, error) {
	beatmap := submission.BeatmapById(beatmapId)
	if beatmap == nil {
		return nil, fmt.Errorf("bss: beatmap not prepared %d", beatmapId)
	}
	if beatmap.Source != nil {
		return nil, fmt.Errorf("bss: beatmap already set: %d", beatmapId)
	}
	if err := submission.CheckBeatmapAccess(beatmap.Target); err != nil {
		return nil, err
	}

	source, err := beatmaps.ParseFromBytes(contents)
	if err != nil {
		return nil, fmt.Errorf("bss: parse beatmap %q: %w", filename, err)
	}

	mode := constants.Mode(source.Mode)
	if !mode.Valid() {
		return nil, fmt.Errorf("bss: invalid beatmap mode %d", source.Mode)
	}

	version := source.Version
	if version == "" {
		version = "Normal"
	}

	target := beatmap.Target
	target.Mode = mode
	target.Version = version
	target.TotalLength = int(source.TotalLength())
	target.DrainLength = int(source.DrainLength())
	target.CountNormal = source.Circles
	target.CountSlider = source.Sliders
	target.CountSpinner = source.Spinners
	target.BPM = source.CommonBPM()
	target.CS = source.Difficulty.GetCS()
	target.AR = source.Difficulty.GetAR()
	target.OD = source.Difficulty.GetOD()
	target.HP = source.Difficulty.GetHP()
	target.SliderMultiplier = source.SliderMultiplier

	source.MapID = int64(beatmap.Target.Id)
	source.SetID = int64(submission.Beatmapset.Id)
	beatmap.Source = source
	beatmap.Target.Filename = filename
	beatmap.Target.Checksum = fmt.Sprintf("%x", md5.Sum(contents))
	return beatmap, nil
}

// CalculateDifficulty calculates modern / eyup star rating
// as well as max combo, and applies the results to the beatmap.
func (p *PreparedBeatmap) CalculateDifficulty(contents []byte, ppv1 *performance.PPv1Service, ppv2 performance.IPPv2Service) error {
	if p.Source == nil {
		return errors.New("bss: beatmap not parsed")
	}

	mode := constants.Mode(p.Source.Mode)
	if !mode.Valid() {
		return fmt.Errorf("bss: invalid beatmap mode %d", p.Source.Mode)
	}
	if ppv1 == nil {
		return errors.New("bss: ppv1 calculator not available")
	}
	if ppv2 == nil || !ppv2.Available() {
		return errors.New("bss: difficulty calculator not available")
	}

	attributes, err := ppv2.CalculateDifficultyFromBytes(
		contents,
		mode,
		constants.NoMod,
	)
	if err != nil {
		return fmt.Errorf("bss: calculate beatmap difficulty: %w", err)
	}

	eyup := ppv1.CalculateEyupStarRating(&schemas.Beatmap{
		Mode:             mode,
		DrainLength:      int(p.Source.DrainLength()),
		CountNormal:      p.Source.Circles,
		CountSlider:      p.Source.Sliders,
		CountSpinner:     p.Source.Spinners,
		BPM:              p.Source.CommonBPM(),
		CS:               p.Source.Difficulty.GetCS(),
		OD:               p.Source.Difficulty.GetOD(),
		HP:               p.Source.Difficulty.GetHP(),
		SliderMultiplier: p.Source.SliderMultiplier,
	})

	p.Target.Diff = attributes.StarRating
	p.Target.DiffEyup = math.Round(eyup*10000) / 10000
	p.Target.MaxCombo = int(attributes.MaxCombo)
	return nil
}
