package bss

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"
)

// SetFS assigns the read-only package used by this submission,
// i.e. either an osz or osz2 archive / filesystem.
func (submission *SubmissionContext) SetFS(source fs.FS) error {
	if submission.FS != nil {
		return errors.New("bss: filesystem already set")
	}
	if submission.Context != nil {
		if err := submission.Context.Err(); err != nil {
			return err
		}
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
