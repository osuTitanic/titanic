package bss

import (
	"archive/zip"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Lekuruu/gosu/pkg/beatmaps"
	"github.com/Lekuruu/osz2-go/pkg/osz2"
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
	"gorm.io/gorm"
)

// ApplyBeatmapsetMetadata applies the generic metadata to the given set.
// Only metadata keys which are present are updated.
func ApplyBeatmapsetMetadata(transaction *state.Repositories, beatmapset *schemas.Beatmapset, metadata osz2.Metadata) error {
	previousName := beatmapset.Name()
	columns := make([]string, 0, 11)
	setMetadata(metadata, osz2.Title, &beatmapset.Title, "title", &columns)
	setMetadata(metadata, osz2.TitleUnicode, &beatmapset.TitleUnicode, "title_unicode", &columns)
	setMetadata(metadata, osz2.Artist, &beatmapset.Artist, "artist", &columns)
	setMetadata(metadata, osz2.ArtistUnicode, &beatmapset.ArtistUnicode, "artist_unicode", &columns)
	setMetadata(metadata, osz2.Creator, &beatmapset.Creator, "creator", &columns)
	setMetadata(metadata, osz2.Source, &beatmapset.Source, "source", &columns)
	setMetadata(metadata, osz2.SourceUnicode, &beatmapset.SourceUnicode, "source_unicode", &columns)
	setMetadata(metadata, osz2.Tags, &beatmapset.Tags, "tags", &columns)

	if len(columns) == 0 {
		// Nothing was updated
		return nil
	}

	nameChanged := beatmapset.Name() != previousName
	if nameChanged {
		// If the artist or title changed, we want to update the display title
		// (https://osu.ppy.sh/wiki/en/Beatmap/Title_text)
		beatmapset.DisplayTitle = new(beatmapset.DefaultDisplayTitle())
		columns = append(columns, "display_title")
	}

	if beatmapset.Tags != nil {
		tags := strings.Split(*beatmapset.Tags, " ")

		detectedLanguage := detectLanguageFromTags(tags)
		detectedGenre := detectGenreFromTags(tags)
		detectedExplicit := detectExplicitFromTags(tags)

		beatmapset.Explicit = beatmapset.Explicit || detectedExplicit
		columns = append(columns, "explicit")

		hasDetectedLanguage := detectedLanguage > constants.BeatmapLanguageUnspecified
		languageAlreadySet := beatmapset.LanguageId > constants.BeatmapLanguageUnspecified

		hasDetectedGenre := detectedGenre > constants.BeatmapGenreUnspecified
		genreAlreadySet := beatmapset.GenreId > constants.BeatmapGenreUnspecified

		if hasDetectedLanguage && !languageAlreadySet {
			beatmapset.LanguageId = detectedLanguage
			columns = append(columns, "language")
		}
		if hasDetectedGenre && !genreAlreadySet {
			beatmapset.GenreId = detectedGenre
			columns = append(columns, "genre")
		}
	}

	rowsAffected, err := transaction.Beatmapsets.Update(beatmapset, columns...)
	return metadataUpdateResult("beatmapset", beatmapset.Id, rowsAffected, err)
}

// ApplySubmissionMetadata applies the metadata that is updated by the
// bss itself rather than the user, e.g. osz sizes, submission status, ...
func ApplySubmissionMetadata(transaction *state.Repositories, submission *SubmissionContext, osz *os.File) (err error) {
	oszInfo, err := osz.Stat()
	if err != nil {
		return fmt.Errorf("inspect osz: %w", err)
	}
	oszSize := oszInfo.Size()
	oszReader, err := zip.NewReader(osz, oszSize)

	submission.Beatmapset.HasVideo, err = CheckFSHasVideo(oszReader)
	if err != nil {
		return err
	}
	submission.Beatmapset.HasStoryboard, err = CheckFSHasOsb(oszReader)
	if err != nil {
		return err
	}

	if !submission.Beatmapset.HasStoryboard {
		// There was no osb included, but maybe the beatmaps contain storyboard events
		submission.Beatmapset.HasStoryboard = submission.BeatmapsHaveSB()
	}

	oszSizeNoVideo := OszSizeWithoutVideo(oszReader, oszSize)
	submission.Beatmapset.OszFilesize = int(oszSize)
	submission.Beatmapset.OszFilesizeNovideo = int(oszSizeNoVideo)

	submission.Beatmapset.LastUpdate = time.Now()
	submission.Beatmapset.Available = true

	previousStatus := submission.Beatmapset.Status
	newStatus := constants.BeatmapStatusWIP

	// NOTE: The user is able to determine if a map is WIP or Pending
	// 		 But this information is only submitted to the forum post
	//       endpoint and not to the upload endpoint itself

	if previousStatus >= constants.BeatmapStatusPending {
		newStatus = constants.BeatmapStatusPending
	}
	submission.Beatmapset.Status = newStatus

	// Apply the new status to every beatmap from the submission as well
	for _, beatmap := range submission.Beatmaps {
		beatmap.Target.Status = newStatus
	}

	rowsAffected, err := transaction.Beatmapsets.Update(submission.Beatmapset,
		"submission_status",
		"has_video",
		"has_storyboard",
		"available",
		"last_updated",
		"osz_filesize",
		"osz_filesize_novideo",
		// TODO: update osz2 hashes maybe
		// (honestly, we don't even use them anymore so who cares)
	)
	return metadataUpdateResult("beatmapset", submission.Beatmapset.Id, rowsAffected, err)
}

// ApplyBeatmapMetadata commits beatmap metadata to the database.
// Please check for beatmap access before calling this.
func ApplyBeatmapMetadata(transaction *state.Repositories, beatmap *schemas.Beatmap) error {
	// NOTE: Beatmap metadata will already be applied
	// 		 by the beatmap parser through `AssignBeatmapContent`
	rowsAffected, err := transaction.Beatmaps.Update(beatmap,
		"status",
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
		"last_updated",
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

func detectExplicitFromTags(tags []string) bool {
	for _, tag := range tags {
		switch strings.ToLower(tag) {
		case "explicit":
			return true
		}
	}
	return false
}

func detectLanguageFromTags(tags []string) constants.BeatmapLanguage {
	for _, tag := range tags {
		switch strings.ToLower(tag) {
		case "english":
			return constants.BeatmapLanguageEnglish
		case "japanese":
			return constants.BeatmapLanguageJapanese
		case "chinese":
			return constants.BeatmapLanguageChinese
		case "instrumental":
			return constants.BeatmapLanguageInstrumental
		case "korean":
			return constants.BeatmapLanguageKorean
		case "french":
			return constants.BeatmapLanguageFrench
		case "german":
			return constants.BeatmapLanguageGerman
		case "swedish":
			return constants.BeatmapLanguageSwedish
		case "italian":
			return constants.BeatmapLanguageItalian
		case "russian":
			return constants.BeatmapLanguageRussian
		case "polish":
			return constants.BeatmapLanguagePolish
		}
	}
	return constants.BeatmapLanguageUnspecified
}

func detectGenreFromTags(tags []string) constants.BeatmapGenre {
	for _, tag := range tags {
		switch strings.ToLower(tag) {
		case "videogame":
			return constants.BeatmapGenreVideoGame
		case "anime":
			return constants.BeatmapGenreAnime
		case "rock":
			return constants.BeatmapGenreRock
		case "pop":
			return constants.BeatmapGenrePop
		case "novelty":
			return constants.BeatmapGenreNovelty
		case "hiphop":
			return constants.BeatmapGenreHipHop
		case "electronic":
			return constants.BeatmapGenreElectronic
		case "metal":
			return constants.BeatmapGenreMetal
		case "classical":
			return constants.BeatmapGenreClassical
		case "folk":
			return constants.BeatmapGenreFolk
		case "jazz":
			return constants.BeatmapGenreJazz
		}
	}
	return constants.BeatmapGenreUnspecified
}
