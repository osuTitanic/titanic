package routes

import (
	"archive/zip"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/state"
	"github.com/osuTitanic/titanic/internal/storage"
	"github.com/osuTitanic/titanic/services/deck/internal/bss"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

type PreOsz2UploadRequest struct {
	Username    string
	Password    string
	Ticket      string
	OszFilename string
	OszTicket   string
	SetId       *int
	First       bool

	File       multipart.File
	FileHeader *multipart.FileHeader
}

func (request *PreOsz2UploadRequest) IsBeatmap() bool {
	if request.FileHeader == nil {
		return false
	}
	fileExtension := path.Ext(request.FileHeader.Filename)
	return strings.EqualFold(fileExtension, ".osu")
}

func (request *PreOsz2UploadRequest) Size() int64 {
	if request.FileHeader == nil {
		return -1
	}
	return request.FileHeader.Size
}

func (request *PreOsz2UploadRequest) Read(limit int64) ([]byte, error) {
	return readFile(request.File, limit)
}

func (request *PreOsz2UploadRequest) Handle(submission *bss.SubmissionContext, store storage.Storage) (*os.File, map[int][]byte, error) {
	if request.IsBeatmap() {
		return handleBeatmapUpload(request, submission, store)
	} else {
		return handleOszUpload(request, submission, store)
	}
}

func (request *PreOsz2UploadRequest) Close() {
	if request == nil {
		return
	}
	if request.File == nil {
		return
	}
	request.File.Close()
	request.File = nil
}

func NewPreOsz2UploadRequest(ctx *server.Context) (request *PreOsz2UploadRequest, err error) {
	err = ctx.EnsureMultipartForm(
		bss.MaxPackageFileSize,
		0, // stream file to disk directly
	)
	if sizeError, ok := errors.AsType[*http.MaxBytesError](err); ok {
		return nil, fmt.Errorf("%w: multipart body exceeds %d bytes", bss.ErrOsz2TooLarge, sizeError.Limit)
	}
	if err != nil {
		return nil, fmt.Errorf("parse multipart form: %w", err)
	}

	request = &PreOsz2UploadRequest{
		Username:    ctx.FormValue("u"),
		Password:    ctx.FormValue("p"),
		Ticket:      ctx.FormValue("c"),
		OszFilename: ctx.FormValue("of"),
		OszTicket:   ctx.FormValue("oc"),
	}
	request.SetId, err = ctx.FormValueIntOptional("s")
	if err != nil {
		return nil, fmt.Errorf("parse beatmapset ID: %w", err)
	}
	if request.SetId != nil && *request.SetId <= 0 {
		return nil, fmt.Errorf("invalid beatmapset ID %d", *request.SetId)
	}

	first, _ := ctx.FormValueIntOptional("r")
	request.First = first != nil && *first == 1

	// Can submit either an .osu file or an .osz file
	// It depends on if this is a full submission or a standard one
	request.File, request.FileHeader, err = ctx.FormFile("osu")
	if err != nil {
		return nil, fmt.Errorf("open upload: %w", err)
	}
	if request.Size() <= 0 {
		size := request.Size()
		request.Close()
		return nil, fmt.Errorf("invalid upload size %d", size)
	}
	if request.Size() > bss.MaxPackageFileSize {
		size := request.Size()
		request.Close()
		return nil, fmt.Errorf("%w: upload is %d bytes", bss.ErrOsz2TooLarge, size)
	}
	return request, nil
}

// /web/osu-bmsubmit-upload.php -> Replaces the complete .osz or inserts one .osu file
func BeatmapSubmissionPreOsz2Upload(ctx *server.Context) {
	request, err := NewPreOsz2UploadRequest(ctx)
	if errors.Is(err, bss.ErrOsz2TooLarge) {
		renderBssErrorText(ctx, err)
		return
	}
	if err != nil {
		renderBssErrorText(ctx, fmt.Errorf("%w: %v", bss.ErrInvalidSubmissionRequest, err))
		return
	}
	defer request.Close()

	user, err := ctx.AuthenticateUser(
		request.Username,
		request.Password,
		false, // does not require a bancho presence
	)
	if errors.Is(err, server.ErrUserNotFound) || errors.Is(err, server.ErrInvalidPassword) {
		renderBssErrorText(ctx, bss.ErrAuthenticationFailed)
		return
	}
	if err != nil {
		renderBssErrorText(ctx, fmt.Errorf("authenticate user: %w", err))
		return
	}

	userPermissions, err := ctx.State.Permissions.Resolve(user.Id)
	if err != nil {
		renderBssErrorText(ctx, fmt.Errorf("resolve permissions for user %d: %w", user.Id, err))
		return
	}

	submission := bss.NewSubmissionContext(ctx.Context())
	submission.User = user

	// Check if user is allowed to upload any beatmaps
	if err := submission.CheckSubmissionEligibility(
		ctx.State.Config.BeatmapSubmissionEnabled,
		userPermissions,
	); err != nil {
		renderBssErrorText(ctx, err)
		return
	}

	setId, err := validateUploadTickets(
		ctx.State.Config.FrontendSecretKey,
		user.Id,
		request,
	)
	if err != nil {
		renderBssErrorText(ctx, err)
		return
	}

	isOszUpload := !request.IsBeatmap()
	newBeatmapset := false

	err = ctx.State.DatabaseTransaction(func(transaction *state.Repositories) error {
		if err := submission.PrepareBeatmapset(transaction, setId); err != nil {
			return err
		}
		if isOszUpload {
			newBeatmapset = submission.Beatmapset.Status == constants.BeatmapStatusInactive
		}

		// Check if we are allowed to update this set
		if err := submission.ResolveAccess(transaction); err != nil {
			return err
		}
		if err := submission.CheckBeatmapsetEligibility(); err != nil {
			return err
		}

		// Create prepared beatmap state from db
		if err := submission.PrepareExistingBeatmaps(); err != nil {
			return err
		}

		// Calls either handleBeatmapUpload or handleOszUpload based on the request type
		// Will produce us an osz file that we can work with to update resources
		oszFile, beatmapFiles, err := request.Handle(submission, ctx.State.Storage)
		if err != nil {
			return err
		}
		if oszFile == nil {
			return nil
		}
		// Cleanup the temporary osz after the transaction
		defer func() {
			oszFile.Close()
			os.Remove(oszFile.Name())
		}()

		info, err := oszFile.Stat()
		if err != nil {
			return fmt.Errorf("bss: inspect legacy osz: %w", err)
		}

		oszSize := info.Size()
		maxBeatmapLength := submission.MaxBeatmapLength()

		// The allowed osz size changes based on the beatmap length
		// The longer the beatmap, the higher the allowed package size
		if err := bss.ValidatePackageSize(oszSize, maxBeatmapLength); err != nil {
			return err
		}

		reader, err := zip.NewReader(oszFile, info.Size())
		if err != nil {
			return fmt.Errorf("bss: open final legacy osz: %w", err)
		}
		if err := bss.ValidateFS(reader); err != nil {
			return fmt.Errorf("bss: validate final legacy osz: %w", err)
		}

		if isOszUpload {
			if err := bss.ApplySubmissionMetadata(transaction, submission, oszFile); err != nil {
				return err
			}
		}

		// Update new .osz file
		if err := submission.SaveOszStream(ctx.State.Storage, oszFile); err != nil {
			return err
		}
		// Upload all new .osu files
		for beatmapId, contents := range beatmapFiles {
			if err := submission.SaveBeatmap(ctx.State.Storage, beatmapId, contents); err != nil {
				return err
			}
		}

		// Generate & upload background thumbnail + audio preview (only if osz updated)
		if isOszUpload && submission.Access.CanUpdateResources {
			if err := submission.UpdateThumbnail(ctx.State.Storage, reader); err != nil {
				ctx.Logger.Warn("Failed to update legacy submission thumbnail", "error", err)
			}
			if err := submission.UpdateAudioPreview(ctx.State.Storage, reader); err != nil {
				ctx.Logger.Warn("Failed to update legacy submission audio preview", "error", err)
			}
		}
		return nil
	})
	if err != nil {
		renderBssErrorText(ctx, err)
		return
	}

	if err := ctx.State.Resources.InvalidateCaches(
		ctx.Context(),
		submission.Beatmapset.Id,
		preparedBeatmapIds(submission)...,
	); err != nil {
		ctx.Logger.Warn(
			"Failed to invalidate resource caches",
			"set_id", submission.Beatmapset.Id, "error", err,
		)
	}

	if isOszUpload || request.First {
		if err := bss.BroadcastActivity(newBeatmapset, ctx.State, submission); err != nil {
			ctx.Logger.Warn(
				"Failed to broadcast beatmap submission activity",
				"user_id", user.Id, "set_id", submission.Beatmapset.Id, "error", err,
			)
		}
	}

	ctx.Logger.Info(
		"Handled pre-osz2 resource upload",
		"user_id", user.Id,
		"set_id", setId,
		"filename", request.FileHeader.Filename,
		"full_package", isOszUpload,
	)
	ctx.RenderText(http.StatusOK, "ok")
}

func handleBeatmapUpload(request *PreOsz2UploadRequest, submission *bss.SubmissionContext, store storage.Storage) (*os.File, map[int][]byte, error) {
	if request.Size() > bss.MaxBeatmapFileSize {
		return nil, nil, fmt.Errorf(
			"%w: beatmap upload is %d bytes",
			bss.ErrOsz2TooLarge, request.Size(),
		)
	}

	beatmap := submission.BeatmapByFilename(request.FileHeader.Filename)
	if beatmap == nil {
		return nil, nil, fmt.Errorf(
			"%w: beatmap %q is not part of set %d",
			bss.ErrInvalidSubmissionRequest, request.FileHeader.Filename, *request.SetId,
		)
	}
	if !submission.Access.CanUpdateBeatmap(beatmap.Target.Id) {
		return nil, nil, nil
	}

	contents, err := request.Read(bss.MaxBeatmapFileSize)
	if err != nil {
		return nil, nil, err
	}
	if !strings.EqualFold(beatmap.Target.Checksum, fmt.Sprintf("%x", md5.Sum(contents))) {
		return nil, nil, fmt.Errorf(
			"%w: beatmap %q does not match prepared metadata",
			bss.ErrInvalidSubmissionRequest, beatmap.Target.Filename,
		)
	}

	oszFile, err := submission.BuildOszBeatmapOverlay(
		store,
		beatmap.Target.Id,
		contents,
		// only remove beatmaps if this is the first submission and we are
		// the set owner. the owner will upload all beatmaps one by one
		request.First && submission.Access.Owner,
	)
	if err != nil {
		return nil, nil, err
	}

	return oszFile, map[int][]byte{beatmap.Target.Id: contents}, nil
}

func handleOszUpload(request *PreOsz2UploadRequest, submission *bss.SubmissionContext, store storage.Storage) (*os.File, map[int][]byte, error) {
	if !strings.EqualFold(path.Ext(request.FileHeader.Filename), ".osz") {
		return nil, nil, fmt.Errorf(
			"%w: unsupported legacy upload %q",
			bss.ErrInvalidSubmissionRequest, request.FileHeader.Filename,
		)
	}

	reader, err := zip.NewReader(request.File, request.Size())
	if err != nil {
		return nil, nil, fmt.Errorf(
			"%w: open osz: %v",
			bss.ErrInvalidSubmissionRequest, err,
		)
	}

	if err := submission.SetFS(reader); err != nil {
		return nil, nil, err
	}
	if err := submission.ValidateFS(); err != nil {
		return nil, nil, fmt.Errorf("%w: validate osz: %v", bss.ErrInvalidSubmissionRequest, err)
	}

	// Read through the osz to gather all beatmaps that we will upload later
	beatmapFiles, err := collectBeatmapsFromOsz(submission, reader)
	if err != nil {
		return nil, nil, err
	}

	// Build the new osz file, accounting for collaboration access
	oszFile, err := submission.BuildOsz(store)
	if err != nil {
		return nil, nil, err
	}
	return oszFile, beatmapFiles, nil
}

func validateUploadTickets(secret string, userId int, request *PreOsz2UploadRequest) (int, error) {
	// The client has received an "upload ticket" in the getid.php request
	// which they can use to "authenticate" / validate their upload request
	//
	// The osz / package ticket is used to identify the target beatmapset
	// The client will use it in every upload.php submission
	// Then there's the individual upload ticket, which identifies either
	// the exact .osu or the complete .osz to upload
	//
	// This ticket system was create presumably such that the transition
	// between getid.php and upload.php is entirely stateless
	packageTicket, err := bss.ParseSubmissionTicket(secret, request.OszTicket)
	if err != nil {
		return 0, fmt.Errorf("validate osz ticket: %w", err)
	}
	if packageTicket.UserId != userId || packageTicket.Kind != bss.SubmissionTicketOsz || packageTicket.Filename != request.OszFilename {
		return 0, bss.ErrInvalidSubmissionTicket
	}
	if request.SetId != nil && *request.SetId != packageTicket.BeatmapsetId {
		return 0, bss.ErrInvalidSubmissionTicket
	}

	uploadTicket, err := bss.ParseSubmissionTicket(secret, request.Ticket)
	if err != nil {
		return 0, fmt.Errorf("validate upload ticket: %w", err)
	}
	expectedKind := bss.SubmissionTicketOsz
	if request.IsBeatmap() {
		expectedKind = bss.SubmissionTicketBeatmap
	}
	if uploadTicket.Kind != expectedKind || uploadTicket.Filename != request.FileHeader.Filename {
		return 0, bss.ErrInvalidSubmissionTicket
	}
	if uploadTicket.UserId != userId || uploadTicket.BeatmapsetId != packageTicket.BeatmapsetId {
		return 0, bss.ErrInvalidSubmissionTicket
	}
	return packageTicket.BeatmapsetId, nil
}

func collectBeatmapsFromOsz(submission *bss.SubmissionContext, source fs.FS) (map[int][]byte, error) {
	contentsById := make(map[int][]byte)

	err := fs.WalkDir(source, ".", func(filename string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filename == "." || entry.IsDir() {
			return nil
		}
		if !strings.EqualFold(path.Ext(filename), ".osu") {
			return nil
		}

		beatmap := submission.BeatmapByFilename(filename)
		if beatmap == nil {
			return nil // probably should return an error here idk
		}
		if !submission.Access.CanUpdateBeatmap(beatmap.Target.Id) {
			return nil
		}

		file, err := source.Open(filename)
		if err != nil {
			return fmt.Errorf("bss: open submitted beatmap %q: %w", filename, err)
		}
		defer file.Close()

		contents, err := readFile(file, bss.MaxBeatmapFileSize)
		if err != nil {
			return fmt.Errorf("bss: read submitted beatmap %q: %w", filename, err)
		}
		if !strings.EqualFold(beatmap.Target.Checksum, fmt.Sprintf("%x", md5.Sum(contents))) {
			return fmt.Errorf("%w: beatmap %q does not match metadata", bss.ErrInvalidSubmissionRequest, filename)
		}
		contentsById[beatmap.Target.Id] = contents
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Check if every submitted map exists in the archive
	for _, beatmap := range submission.Beatmaps {
		if !submission.Access.CanUpdateBeatmap(beatmap.Target.Id) {
			continue
		}
		if _, exists := contentsById[beatmap.Target.Id]; !exists {
			return nil, fmt.Errorf("%w: missing beatmap %q", bss.ErrInvalidSubmissionRequest, beatmap.Target.Filename)
		}
	}
	return contentsById, nil
}

func readFile(source io.Reader, limit int64) ([]byte, error) {
	contents, err := io.ReadAll(io.LimitReader(source, limit+1))
	if err != nil {
		return nil, fmt.Errorf("bss: read legacy upload: %w", err)
	}
	if int64(len(contents)) > limit {
		return nil, bss.ErrOsz2TooLarge
	}
	if len(contents) == 0 {
		return nil, fmt.Errorf("%w: uploaded file is empty", bss.ErrInvalidSubmissionRequest)
	}
	return contents, nil
}
