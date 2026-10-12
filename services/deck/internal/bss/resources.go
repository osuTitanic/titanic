package bss

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
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

// PutFSWithAccess writes the package contents accounting for
// collaboration status, which determines what they can edit
func (submission *SubmissionContext) PutFSWithAccess(writer *zip.Writer, existingFS fs.FS) error {
	if submission.Access == nil {
		return errors.New("bss: submission access not resolved")
	}
	if submission.Access.Owner {
		// We are the owner of this mapset, we can do anything™
		return submission.PutFS(writer, nil)
	}

	// We are a collaborator, and are limited to what we are
	// able to update in this beatmapset

	retainedBeatmaps := make(map[string]bool)
	submittedBeatmaps := make(map[string]bool)

	for _, beatmap := range submission.Beatmaps {
		filename := beatmap.Target.Filename
		if submission.Access.CanUpdateBeatmap(beatmap.Target.Id) {
			// Collaborator is allowed to update this beatmap
			submittedBeatmaps[filename] = true
			continue
		}

		info, err := fs.Stat(existingFS, filename)
		if err != nil || !info.Mode().IsRegular() {
			return fmt.Errorf("bss: existing package has no beatmap file for ID %d", beatmap.Target.Id)
		}

		// Retain the original file, because collaborator is not allowed to update it
		retainedBeatmaps[filename] = true
	}

	putFiles := func(source fs.FS, beatmaps map[string]bool, includeResources bool) error {
		return submission.PutFSFrom(writer, source, func(filename string) bool {
			if strings.EqualFold(path.Ext(filename), ".osu") {
				// Can we update this beatmap?
				return beatmaps[filename]
			}
			// Are we allowed to update resources (e.g. background / video / sb assets)
			return includeResources
		})
	}

	// Write all beatmaps that should be retained
	// If the collaborator can update resources: don't include the existing resources
	// If the collaborator can't update resources: include the existing resources
	if err := putFiles(existingFS, retainedBeatmaps, !submission.Access.CanUpdateResources); err != nil {
		return err
	}

	// Finally, write the beatmaps that the collaborator is able to update
	// as well as resource files, if they are allowed to update those
	return putFiles(submission.FS, submittedBeatmaps, submission.Access.CanUpdateResources)
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

// BuildOsz writes the submitted filesystem into a temporary osz.
// The returned file must be closed & removed manually.
func (submission *SubmissionContext) BuildOsz(store storage.Storage) (output *os.File, err error) {
	if submission.Access == nil {
		return nil, errors.New("bss: submission access not resolved")
	}
	if submission.Beatmapset == nil {
		return nil, ErrBeatmapsetNotPrepared
	}
	var existingFS fs.FS

	if !submission.Access.Owner {
		// We are dealing with a collaborator of this set
		// They can't replace the entire osz file so we
		// need to use the existing osz as a base
		existingSource, existingSize, err := store.ReadStreamAt(
			submission.Context,
			strconv.Itoa(submission.Beatmapset.Id),
			"osz",
		)
		if err != nil {
			return nil, fmt.Errorf("bss: read existing osz: %w", err)
		}
		defer existingSource.Close()

		existingReader, err := zip.NewReader(existingSource, existingSize)
		if err != nil {
			return nil, fmt.Errorf("bss: open existing osz: %w", err)
		}
		existingFS = existingReader
	}

	// Create the temporary osz file handle
	// (the osz will be uploaded to persistent storage later)
	temporary, err := os.CreateTemp("", "bss-osz-*")
	if err != nil {
		return nil, fmt.Errorf("bss: create temporary osz file: %w", err)
	}
	defer func() {
		// Only close & remove the temporary osz if we got an error
		// Otherwise we want to keep it open for processing
		if err != nil {
			err = errors.Join(err, temporary.Close(), os.Remove(temporary.Name()))
		}
	}()

	// Write the new osz file
	writer := zip.NewWriter(temporary)

	if writeErr := submission.PutFSWithAccess(writer, existingFS); writeErr != nil {
		err = fmt.Errorf("bss: write osz: %w", errors.Join(writeErr, writer.Close()))
		return nil, err
	}
	if err = writer.Close(); err != nil {
		return nil, fmt.Errorf("bss: close osz: %w", err)
	}
	if _, err = temporary.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("bss: rewind osz file: %w", err)
	}
	return temporary, nil
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

// CheckFSIncludes checks whether the fs contains a file specified by include.
func CheckFSIncludes(source fs.FS, include func(filename string) bool) (bool, error) {
	if source == nil {
		return false, errors.New("bss: filesystem not set")
	}
	found := false

	err := fs.WalkDir(source, ".", func(filename string, entry fs.DirEntry, err error) error {
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

// CheckFSHasVideo reports whether the fs contains a video.
func CheckFSHasVideo(source fs.FS) (bool, error) {
	return CheckFSIncludes(source, resources.IsVideoFile)
}

// CheckFSHasOsb checks whether the fs contains a storyboard file.
func CheckFSHasOsb(source fs.FS) (bool, error) {
	return CheckFSIncludes(source, resources.IsStoryboardFile)
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
