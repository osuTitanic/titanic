package bss

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strconv"
	"strings"

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
func (submission *SubmissionContext) PutFS(writer *zip.Writer, include func(filename string) bool) error {
	if submission.FS == nil {
		return errors.New("bss: filesystem not set")
	}
	if submission.IsCanceled() {
		return submission.Context.Err()
	}

	return fs.WalkDir(submission.FS, ".", func(filename string, entry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("bss: walk package file %q: %w", filename, err)
		}
		if filename == "." || entry.IsDir() {
			return nil
		}
		if submission.IsCanceled() {
			return submission.Context.Err()
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

		target, err := writer.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("bss: create package file %q: %w", filename, err)
		}
		source, err := submission.FS.Open(filename)
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
	if err := submission.validateResourceStorage(store); err != nil {
		return err
	}

	if err := store.Save(
		submission.Context,
		strconv.Itoa(submission.Beatmapset.Id),
		"osz",
		contents,
	); err != nil {
		return fmt.Errorf("bss: save osz for set %d: %w", submission.Beatmapset.Id, err)
	}
	return nil
}

// SaveOsz2 creates / overwrites the osz2 package for the prepared set.
func (submission *SubmissionContext) SaveOsz2(store storage.Storage, contents []byte) error {
	if err := submission.validateResourceStorage(store); err != nil {
		return err
	}

	if err := store.Save(
		submission.Context,
		strconv.Itoa(submission.Beatmapset.Id),
		"osz2",
		contents,
	); err != nil {
		return fmt.Errorf("bss: save osz2 for set %d: %w", submission.Beatmapset.Id, err)
	}
	return nil
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
