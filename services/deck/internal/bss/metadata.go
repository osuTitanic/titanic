package bss

import (
	"fmt"
	"strconv"

	"github.com/Lekuruu/gosu/pkg/beatmaps"
	"github.com/Lekuruu/osz2-go/pkg/osz2"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
	"gorm.io/gorm"
)

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

// ApplyBeatmapMetadata commits beatmap metadata to the database.
// Please check for beatmap access before calling this.
func ApplyBeatmapMetadata(repositories *state.Repositories, beatmap *schemas.Beatmap) error {
	// NOTE: Beatmap metadata will already be applied
	// 		 by the beatmap parser through `AssignBeatmapContent`
	rowsAffected, err := repositories.Beatmaps.Update(beatmap,
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
	return metadataUpdateResult("beatmap", beatmap.Id, rowsAffected, err)
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
