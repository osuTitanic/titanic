package bss

import (
	"errors"
	"fmt"
	"io/fs"
	"math"
	"path"
	"strings"

	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
)

const (
	MaxBeatmapFileSize = 15_000_000
	MaxPackageFileSize = 100_000_000
)

var supportedPackageExtensions = map[string]bool{
	".3gp":  true,
	".avi":  true,
	".flac": true,
	".flv":  true,
	".ini":  true,
	".jpeg": true,
	".jpg":  true,
	".m4v":  true,
	".mkv":  true,
	".mov":  true,
	".mp3":  true,
	".mp4":  true,
	".mpeg": true,
	".mpg":  true,
	".ogg":  true,
	".ogv":  true,
	".osb":  true,
	".osk":  true,
	".osu":  true,
	".osz":  true,
	".png":  true,
	".wav":  true,
	".webm": true,
	".wmv":  true,
}

// ValidatePackageSize validates a serialized .osz
// against the duration-based package size limit.
func ValidatePackageSize(packageSize int64, maxBeatmapLength float64) error {
	if maxBeatmapLength <= 1 || math.IsNaN(maxBeatmapLength) {
		return ErrBeatmapTooShort
	}
	if packageSize < 0 || packageSize > PackageSizeLimit(maxBeatmapLength) {
		return ErrPackageTooLarge
	}
	return nil
}

// PackageSizeLimit returns the compressed .osz size limit based on the
// longest difficulty in the set. Beatmap length is given in seconds.
func PackageSizeLimit(maxBeatmapLength float64) int64 {
	if maxBeatmapLength <= 0 || math.IsNaN(maxBeatmapLength) {
		return 10_000_000
	}
	if math.IsInf(maxBeatmapLength, 1) {
		return MaxPackageFileSize
	}

	limit := 15_000_000 + int64(10_000_000*(maxBeatmapLength/60))
	return min(limit, int64(MaxPackageFileSize))
}

// ValidateBeatmapUniqueness rejects submitted filenames / checksums already
// belonging to another mapper, while allowing the same mapper to re-use them.
func (submission *SubmissionContext) ValidateBeatmapUniqueness(repositories *state.Repositories) error {
	if submission.Beatmapset == nil {
		return ErrBeatmapsetNotPrepared
	}

	isDuplicateBeatmap := func(submittedId int, existing *schemas.Beatmap) bool {
		if existing == nil {
			return false
		}
		if existing.Id == submittedId {
			return false
		}
		if existing.SetId == submission.Beatmapset.Id {
			return true
		}
		if existing.Beatmapset.CreatorId == nil || submission.Beatmapset.CreatorId == nil {
			// should not happen... well... *should*
			return true
		}
		return *existing.Beatmapset.CreatorId != *submission.Beatmapset.CreatorId
	}

	for _, beatmap := range submission.Beatmaps {
		if beatmap.Source == nil {
			continue
		}

		// Check if any other beatmap in this submission has the same filename
		for _, other := range submission.Beatmaps {
			if other.Target.Id == beatmap.Target.Id {
				continue
			}
			if strings.EqualFold(other.Target.Filename, beatmap.Target.Filename) {
				return ErrDuplicateBeatmap
			}
		}

		byFilename, err := repositories.Beatmaps.ByFilename(beatmap.Target.Filename, "Beatmapset")
		if err != nil {
			return fmt.Errorf("bss: beatmap by filename %q: %w", beatmap.Target.Filename, err)
		}

		// Check if the filename belongs to another mapper's beatmap,
		// or if it belongs to a different beatmap in the same set
		hasFilenameConflict := submission.checkFilenameConflicts(beatmap.Target, byFilename)
		duplicateBeatmap := isDuplicateBeatmap(beatmap.Target.Id, byFilename)

		if hasFilenameConflict && duplicateBeatmap {
			return ErrDuplicateBeatmap
		}

		byChecksum, err := repositories.Beatmaps.ByChecksum(beatmap.Target.Checksum, "Beatmapset")
		if err != nil {
			return fmt.Errorf("bss: beatmap by checksum %q: %w", beatmap.Target.Checksum, err)
		}
		if isDuplicateBeatmap(beatmap.Target.Id, byChecksum) {
			return ErrDuplicateBeatmap
		}
	}
	return nil
}

func (submission *SubmissionContext) checkFilenameConflicts(submitted, existing *schemas.Beatmap) bool {
	if existing == nil {
		return false
	}
	if existing.Id == submitted.Id {
		return false
	}
	if existing.SetId == submission.Beatmapset.Id {
		// If the existing beatmap belongs to the same set, we allow the filename to be reused
		if prepared := submission.BeatmapById(existing.Id); prepared != nil && prepared.Source != nil {
			existing = prepared.Target
		}
	}
	return strings.EqualFold(submitted.Filename, existing.Filename)
}

// ValidateFS validates all files in the assigned package.
func (submission *SubmissionContext) ValidateFS() error {
	if submission.FS == nil {
		return errors.New("bss: package filesystem not set")
	}
	if submission.IsCanceled() {
		return submission.Context.Err()
	}
	return ValidateFS(submission.FS)
}

// ValidateFS validates all files in a read-only beatmap package.
func ValidateFS(source fs.FS) error {
	paths := make(map[string]packagePath)
	hasBeatmap := false
	files := 0

	err := fs.WalkDir(source, ".", func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("bss: walk package file %q: %w", name, err)
		}
		if name == "." {
			return nil
		}

		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("bss: inspect package file %q: %w", name, err)
		}

		isDir := info.IsDir()
		name = strings.TrimSuffix(name, "/")

		// Check name for invalid characters and path conflicts
		if err := validatePackageFilename(name); err != nil {
			return fmt.Errorf("bss: invalid package file: %w", err)
		}

		// Track paths & detect path conflicts
		if err := addPackagePath(paths, name, isDir); err != nil {
			return err
		}
		if isDir {
			return nil
		}

		// non-regular e.g. symlinks and device files should not be present
		if !info.Mode().IsRegular() || info.Size() < 0 {
			return fmt.Errorf("bss: invalid package file: %q is not a regular file", name)
		}
		extension := strings.ToLower(path.Ext(name))

		if !supportedPackageExtensions[extension] {
			return fmt.Errorf("bss: unsupported package file: %q", name)
		}

		files++
		hasBeatmap = hasBeatmap || extension == ".osu"
		return nil
	})
	if err != nil {
		return err
	}

	if files == 0 {
		return errors.New("bss: package is empty")
	}
	if !hasBeatmap {
		return errors.New("bss: package has no beatmaps")
	}
	return nil
}

type packagePath struct {
	name  string
	isDir bool
}

func addPackagePath(paths map[string]packagePath, name string, isDir bool) error {
	for current := name; current != "."; current = path.Dir(current) {
		currentIsDir := current != name || isDir
		key := strings.ToLower(current)

		existing, exists := paths[key]
		if !exists {
			// first occurrence of this path
			paths[key] = packagePath{
				name:  current,
				isDir: currentIsDir,
			}
			continue
		}

		// there should not be a name or isDir mismatch
		// e.g. "Maps/beatmap.osu" and "maps/audio.mp3"
		// or "audio.mp3" and "audio.mp3/file.osu"
		if existing.name != current || existing.isDir != currentIsDir || !currentIsDir {
			return fmt.Errorf(
				"bss: duplicate package file: %q conflicts with %q",
				name,
				existing.name,
			)
		}
	}

	return nil
}

func validatePackageFilename(filename string) error {
	if filename == "" || filename == "." || strings.HasSuffix(filename, "/") {
		return fmt.Errorf("bss: invalid package filename %q", filename)
	}
	if path.IsAbs(filename) || path.Clean(filename) != filename || !fs.ValidPath(filename) {
		return fmt.Errorf("bss: invalid package filename %q", filename)
	}
	if strings.ContainsAny(filename, "\\\x00") {
		return fmt.Errorf("bss: invalid package filename %q", filename)
	}
	if first, _, _ := strings.Cut(filename, "/"); strings.Contains(first, ":") {
		return fmt.Errorf("bss: invalid package filename %q", filename)
	}
	return nil
}
