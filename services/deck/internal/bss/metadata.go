package bss

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Lekuruu/gosu/pkg/beatmaps"
	"github.com/Lekuruu/osz2-go/pkg/osz2"
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
	"gorm.io/gorm"
)

// TODO: Not sure yet if we want to have a
// 		 struct like this when we could derive
// 		 metadata in ApplyBeatmapMetadata

type BeatmapMetadata struct {
	Filename       string
	Checksum       string
	BPM            float64
	TotalLength    int
	DrainLength    int
	MaxCombo       int
	Difficulty     float64
	DifficultyEyup float64
}

// SetBeatmap assigns a parsed beatmap & its metadata to a beatmapId.
func SetBeatmap(submission *SubmissionContext, beatmapId int, source *beatmaps.Beatmap, metadata BeatmapMetadata) (*PreparedBeatmap, error) {
	for _, beatmap := range submission.Beatmaps {
		if beatmap.Target.Id != beatmapId {
			continue
		}
		if err := CheckBeatmapAccess(submission, beatmap.Target); err != nil {
			return nil, err
		}
		if beatmap.Source != nil {
			return nil, fmt.Errorf("bss: beatmap already set: %d", beatmapId)
		}

		source.MapID = int64(beatmap.Target.Id)
		source.SetID = int64(submission.Beatmapset.Id)
		beatmap.Source = source
		beatmap.Metadata = metadata
		return beatmap, nil
	}
	return nil, fmt.Errorf("bss: beatmap not prepared %d", beatmapId)
}

// ApplyBeatmapsetMetadata applies the metadata to the given set.
// Only metadata keys which are present are updated.
func ApplyBeatmapsetMetadata(repositories *state.Repositories, beatmapset *schemas.Beatmapset, metadata osz2.Metadata) error {
	columns := make([]string, 0, 8)
	setMetadata(metadata, osz2.Title, &beatmapset.Title, "title", &columns)
	setMetadata(metadata, osz2.TitleUnicode, &beatmapset.TitleUnicode, "title_unicode", &columns)
	setMetadata(metadata, osz2.Artist, &beatmapset.Artist, "artist", &columns)
	setMetadata(metadata, osz2.ArtistUnicode, &beatmapset.ArtistUnicode, "artist_unicode", &columns)
	setMetadata(metadata, osz2.Creator, &beatmapset.Creator, "creator", &columns)
	setMetadata(metadata, osz2.Source, &beatmapset.Source, "source", &columns)
	setMetadata(metadata, osz2.SourceUnicode, &beatmapset.SourceUnicode, "source_unicode", &columns)
	setMetadata(metadata, osz2.Tags, &beatmapset.Tags, "tags", &columns)

	if len(columns) == 0 {
		return nil
	}
	rowsAffected, err := repositories.Beatmapsets.Update(beatmapset, columns...)
	return metadataUpdateResult("beatmapset", beatmapset.Id, rowsAffected, err)
}

// ApplyBeatmapMetadata applies a prepared beatmap's parsed & derived metadata.
func ApplyBeatmapMetadata(submission *SubmissionContext, repositories *state.Repositories, beatmap *PreparedBeatmap) error {
	if beatmap.Source == nil {
		return errors.New("bss: beatmap not set")
	}
	if err := CheckBeatmapAccess(submission, beatmap.Target); err != nil {
		return err
	}
	target := beatmap.Target
	source := beatmap.Source
	metadata := beatmap.Metadata

	mode := constants.Mode(source.Mode)
	if !mode.Valid() {
		return fmt.Errorf("bss: invalid beatmap mode %d", source.Mode)
	}

	version := source.Version
	if version == "" {
		version = "Normal"
	}

	target.Mode = mode
	target.Checksum = metadata.Checksum
	target.Version = version
	target.Filename = metadata.Filename
	target.TotalLength = metadata.TotalLength
	target.DrainLength = metadata.DrainLength
	target.CountNormal = source.Circles
	target.CountSlider = source.Sliders
	target.CountSpinner = source.Spinners
	target.MaxCombo = metadata.MaxCombo
	target.BPM = metadata.BPM
	target.CS = source.Difficulty.GetCS()
	target.AR = source.Difficulty.GetAR()
	target.OD = source.Difficulty.GetOD()
	target.HP = source.Difficulty.GetHP()
	target.Diff = metadata.Difficulty
	target.DiffEyup = metadata.DifficultyEyup
	target.SliderMultiplier = source.SliderMultiplier

	rowsAffected, err := repositories.Beatmaps.Update(target,
		"mode",
		"md5",
		"version",
		"filename",
		"total_length",
		"drain_length",
		"count_normal",
		"count_slider",
		"count_spinner",
		"max_combo",
		"bpm",
		"cs",
		"ar",
		"od",
		"hp",
		"diff",
		"diff_eyup",
		"slider_multiplier",
	)
	return metadataUpdateResult("beatmap", target.Id, rowsAffected, err)
}

// MetadataFromBeatmap converts the set metadata embedded in a .osu
// file into the same representation used by osz2 packages.
func MetadataFromBeatmap(beatmap *beatmaps.Beatmap) osz2.Metadata {
	return osz2.Metadata{
		osz2.Title:         beatmap.Title,
		osz2.TitleUnicode:  beatmap.TitleUnicode,
		osz2.Artist:        beatmap.Artist,
		osz2.ArtistUnicode: beatmap.ArtistUnicode,
		osz2.Creator:       beatmap.Creator,
		osz2.Source:        beatmap.Source,
		osz2.Tags:          beatmap.Tags,
		osz2.BeatmapSetID:  strconv.FormatInt(beatmap.SetID, 10),
	}
}

func setMetadata(metadata osz2.Metadata, metadataType osz2.MetaType, target **string, column string, columns *[]string) {
	value, ok := metadata[metadataType]
	if ok {
		*target = &value
		*columns = append(*columns, column)
	}
}

func metadataUpdateResult(schema string, id int, rowsAffected int64, err error) error {
	if err != nil {
		return fmt.Errorf("bss: update %s %d metadata: %w", schema, id, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("bss: update %s %d metadata: %w", schema, id, gorm.ErrRecordNotFound)
	}
	if rowsAffected != 1 {
		return fmt.Errorf("bss: update %s %d metadata: affected %d rows", schema, id, rowsAffected)
	}
	return nil
}
