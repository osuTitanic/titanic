package bss

import (
	"errors"
	"fmt"
	"io/fs"
	"math"
	"path"
	"strings"
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
		return errors.New("bss: beatmap is too short")
	}
	if packageSize < 0 || packageSize > PackageSizeLimit(maxBeatmapLength) {
		return errors.New("bss: package is too large")
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

	limit := 10_000_000 + int64(10_000_000*(maxBeatmapLength/60))
	return min(limit, int64(MaxPackageFileSize))
}

// ValidateFS validates all files in the assigned package.
func (submission *SubmissionContext) ValidateFS() error {
	if submission.FS == nil {
		return errors.New("bss: package filesystem not set")
	}
	if submission.Context != nil {
		if err := submission.Context.Err(); err != nil {
			return err
		}
	}

	paths := make(map[string]packagePath)
	hasBeatmap := false
	files := 0

	err := fs.WalkDir(submission.FS, ".", func(name string, entry fs.DirEntry, err error) error {
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
