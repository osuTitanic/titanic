package bss

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strconv"
	"strings"

	"github.com/Lekuruu/zipstream"
	"github.com/osuTitanic/titanic/internal/resources"
	"github.com/osuTitanic/titanic/internal/storage"
)

// SetFS assigns the read-only package used by this submission,
// i.e. either an osz or osz2 archive / filesystem.
func (submission *SubmissionContext) SetFS(source fs.FS) error {
	if submission.FS != nil {
		return errors.New("bss: filesystem already set")
	}
	if submission.IsCanceled() {
		return submission.Context.Err()
	}

	root, err := fs.Stat(source, ".")
	if err != nil {
		return fmt.Errorf("bss: open package filesystem: %w", err)
	}
	if !root.IsDir() {
		return errors.New("bss: package filesystem root is not a directory")
	}

	submission.FS = source
	return nil
}

// PutFS copies included files from the submission filesystem into an osz.
func (submission *SubmissionContext) PutFS(target *zip.Writer, include func(filename string) bool) error {
	return submission.PutFSFrom(target, submission.FS, include)
}

// PutFSFrom copies included files from a read-only filesystem into an osz.
func (submission *SubmissionContext) PutFSFrom(target *zip.Writer, sourceFS fs.FS, include func(filename string) bool) error {
	if sourceFS == nil {
		return errors.New("bss: filesystem not set")
	}
	if submission.IsCanceled() {
		return submission.Context.Err()
	}

	return fs.WalkDir(sourceFS, ".", func(filename string, entry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("bss: walk package file %q: %w", filename, err)
		}
		if filename == "." || entry.IsDir() {
			return nil
		}

		if include != nil && !include(filename) {
			return nil
		}
		if err := validatePackageFilename(filename); err != nil {
			return err
		}

		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("bss: get package file info %q: %w", filename, err)
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("bss: package file %q is not regular file", filename)
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("bss: create package header %q: %w", filename, err)
		}
		header.Name = filename
		header.Method = zip.Deflate

		target, err := target.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("bss: create package file %q: %w", filename, err)
		}
		source, err := sourceFS.Open(filename)
		if err != nil {
			return fmt.Errorf("bss: open package file %q: %w", filename, err)
		}

		written, copyErr := io.Copy(target, source)
		closeErr := source.Close()
		if err := errors.Join(copyErr, closeErr); err != nil {
			return fmt.Errorf("bss: write package file %q: %w", filename, err)
		}
		if written != info.Size() {
			return fmt.Errorf("bss: write package file %q: %w", filename, io.ErrUnexpectedEOF)
		}
		return nil
	})
}

// CheckFSIncludes checks whether the submission contains a file specified by include.
func (submission *SubmissionContext) CheckFSIncludes(include func(filename string) bool) (bool, error) {
	if submission.FS == nil {
		return false, errors.New("bss: filesystem not set")
	}
	if submission.IsCanceled() {
		return false, submission.Context.Err()
	}
	found := false

	err := fs.WalkDir(submission.FS, ".", func(filename string, entry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("bss: walk package file %q: %w", filename, err)
		}
		if filename == "." || entry.IsDir() {
			return nil
		}

		if include(filename) {
			found = true
			return fs.SkipAll
		}
		return nil
	})
	return found, err
}

// CheckFSHasVideo reports whether the submission contains a video.
func (submission *SubmissionContext) CheckFSHasVideo() (bool, error) {
	return submission.CheckFSIncludes(resources.IsVideoFile)
}

// CheckFSHasOsb checks whether the submission contains a storyboard file.
func (submission *SubmissionContext) CheckFSHasOsb() (bool, error) {
	return submission.CheckFSIncludes(resources.IsStoryboardFile)
}

// PutBeatmap writes a prepared difficulty to an
// osz archive without changing its metadata.
func (submission *SubmissionContext) PutBeatmap(writer *zip.Writer, beatmapId int, contents []byte) error {
	beatmap := submission.BeatmapById(beatmapId)
	if beatmap == nil {
		return fmt.Errorf("bss: beatmap not prepared %d", beatmapId)
	}
	if err := submission.CheckBeatmapAccess(beatmap.Target); err != nil {
		return err
	}
	filename := beatmap.Target.Filename

	if !strings.EqualFold(path.Ext(filename), ".osu") {
		return fmt.Errorf("bss: invalid beatmap filename %q", filename)
	}
	return PutFile(writer, filename, contents)
}

// PutFile writes one regular file to an osz archive.
func PutFile(writer *zip.Writer, filename string, contents []byte) error {
	if err := validatePackageFilename(filename); err != nil {
		return err
	}

	entry, err := writer.Create(filename)
	if err != nil {
		return fmt.Errorf("bss: create package file %q: %w", filename, err)
	}
	written, err := entry.Write(contents)
	if err != nil {
		return fmt.Errorf("bss: write package file %q: %w", filename, err)
	}
	if written != len(contents) {
		return fmt.Errorf("bss: write package file %q: %w", filename, io.ErrShortWrite)
	}
	return nil
}

// SaveOsz creates / overwrites the osz package for the prepared beatmapset.
func (submission *SubmissionContext) SaveOsz(store storage.Storage, contents []byte) error {
	return submission.SaveOszStream(store, bytes.NewReader(contents))
}

// SaveOszStream creates / overwrites the osz package for the prepared beatmapset.
func (submission *SubmissionContext) SaveOszStream(store storage.Storage, stream io.Reader) error {
	if err := submission.validateResourceStorage(store); err != nil {
		return err
	}
	return store.SaveStream(
		submission.Context,
		strconv.Itoa(submission.Beatmapset.Id),
		"osz",
		stream,
	)
}

// SaveOsz2 creates / overwrites the osz2 package for the prepared set.
func (submission *SubmissionContext) SaveOsz2(store storage.Storage, contents []byte) error {
	return submission.SaveOsz2Stream(store, bytes.NewReader(contents))
}

// SaveOsz2Stream creates / overwrites the osz2 package for the prepared set.
func (submission *SubmissionContext) SaveOsz2Stream(store storage.Storage, stream io.Reader) error {
	if err := submission.validateResourceStorage(store); err != nil {
		return err
	}
	return store.SaveStream(
		submission.Context,
		strconv.Itoa(submission.Beatmapset.Id),
		"osz2",
		stream,
	)
}

// SaveBeatmap overwrites a .osu file without changing its metadata.
func (submission *SubmissionContext) SaveBeatmap(store storage.Storage, beatmapId int, contents []byte) error {
	if err := submission.validateResourceStorage(store); err != nil {
		return err
	}

	beatmap := submission.BeatmapById(beatmapId)
	if beatmap == nil {
		return fmt.Errorf("bss: beatmap not prepared %d", beatmapId)
	}
	if err := submission.CheckBeatmapAccess(beatmap.Target); err != nil {
		return err
	}

	if err := store.Save(
		submission.Context,
		strconv.Itoa(beatmap.Target.Id),
		"beatmaps",
		contents,
	); err != nil {
		return fmt.Errorf("bss: save beatmap %d: %w", beatmap.Target.Id, err)
	}
	return nil
}

// OszSizeWithoutVideo returns the size of a zip archive with videos removed.
func OszSizeWithoutVideo(reader *zip.Reader, fallbackSize int64) int64 {
	size, ok := zipstream.ArchiveSizeFiltered(reader.File, func(file *zip.File) bool {
		return !resources.IsVideoFile(file.Name)
	})
	if !ok {
		return fallbackSize
	}
	return size
}

func (submission *SubmissionContext) validateResourceStorage(store storage.Storage) error {
	if store == nil {
		return errors.New("bss: nil resource storage")
	}
	if submission.Context == nil {
		return errors.New("bss: nil submission context")
	}
	if submission.IsCanceled() {
		return submission.Context.Err()
	}
	if submission.Beatmapset == nil || submission.Beatmapset.Id <= 0 {
		return errors.New("bss: beatmapset not prepared")
	}
	return nil
}

func (submission *SubmissionContext) readFile(filename string, maxSize int64) (contents []byte, err error) {
	file, err := submission.FS.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = errors.Join(err, file.Close())
	}()

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() || info.Size() < 0 {
		return nil, errors.New("not a regular file")
	}
	if info.Size() > maxSize {
		return nil, fmt.Errorf("file exceeds maximum size of %d bytes", maxSize)
	}

	contents, err = io.ReadAll(io.LimitReader(file, maxSize+1))
	if err != nil {
		return nil, err
	}
	if int64(len(contents)) > maxSize {
		return nil, fmt.Errorf("file exceeds maximum size of %d bytes", maxSize)
	}
	return contents, nil
}
