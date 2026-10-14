package routes

import (
	"archive/zip"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/osuTitanic/titanic/internal/state"
	"github.com/osuTitanic/titanic/internal/storage"
	"github.com/osuTitanic/titanic/services/deck/internal/bss"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

type UploadHandler func(
	request *PreOsz2UploadRequest,
	submission *bss.SubmissionContext,
	store storage.Storage,
) (*os.File, error)

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
	contents, err := io.ReadAll(io.LimitReader(request.File, limit+1))
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

func (request *PreOsz2UploadRequest) Handle(submission *bss.SubmissionContext, store storage.Storage) (*os.File, error) {
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

	err = ctx.State.DatabaseTransaction(func(transaction *state.Repositories) error {
		if err := submission.PrepareBeatmapset(transaction, setId); err != nil {
			return err
		}
		if err := submission.ResolveAccess(transaction); err != nil {
			return err
		}
		if err := submission.CheckBeatmapsetEligibility(); err != nil {
			return err
		}
		if err := submission.PrepareExistingBeatmaps(); err != nil {
			return err
		}

		// Calls either handleBeatmapUpload or handleOszUpload based on the request type
		// Will produce us an osz file that we can work with to update resources
		oszFile, err := request.Handle(submission, ctx.State.Storage)
		if err != nil {
			return err
		}
		defer func() {
			oszFile.Close()
			os.Remove(oszFile.Name())
		}()

		info, err := oszFile.Stat()
		if err != nil {
			return fmt.Errorf("bss: inspect legacy osz: %w", err)
		}
		if err := bss.ValidatePackageSize(info.Size(), submission.MaxBeatmapLength()); err != nil {
			return err
		}

		reader, err := zip.NewReader(oszFile, info.Size())
		if err != nil {
			return fmt.Errorf("bss: open final legacy osz: %w", err)
		}
		if err := bss.ValidateFS(reader); err != nil {
			return fmt.Errorf("bss: validate final legacy osz: %w", err)
		}

		if err := submission.SaveOszStream(ctx.State.Storage, oszFile); err != nil {
			return err
		}

		if !request.IsBeatmap() && submission.Access.CanUpdateResources {
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

	// TODO: Invalidate cache

	ctx.Logger.Info(
		"Handled pre-osz2 resource upload",
		"user_id", user.Id,
		"set_id", setId,
		"filename", request.FileHeader.Filename,
		"full_package", !request.IsBeatmap(),
	)
	ctx.RenderText(http.StatusOK, "ok")
}

func handleBeatmapUpload(request *PreOsz2UploadRequest, submission *bss.SubmissionContext, store storage.Storage) (*os.File, error) {
	if request.Size() > bss.MaxBeatmapFileSize {
		return nil, fmt.Errorf("%w: beatmap upload is %d bytes", bss.ErrOsz2TooLarge, request.Size())
	}

	beatmap := submission.BeatmapByFilename(request.FileHeader.Filename)
	if beatmap == nil {
		return nil, fmt.Errorf(
			"%w: beatmap %q is not part of set %d",
			bss.ErrInvalidSubmissionRequest, request.FileHeader.Filename, *request.SetId,
		)
	}
	if !submission.Access.CanUpdateBeatmap(beatmap.Target.Id) {
		return nil, nil
	}

	contents, err := request.Read(bss.MaxBeatmapFileSize)
	if err != nil {
		return nil, err
	}
	if !strings.EqualFold(beatmap.Target.Checksum, fmt.Sprintf("%x", md5.Sum(contents))) {
		return nil, fmt.Errorf("%w: beatmap %q does not match prepared metadata", bss.ErrInvalidSubmissionRequest, beatmap.Target.Filename)
	}

	oszFile, err := submission.BuildOszBeatmapOverlay(
		store,
		beatmap.Target.Id,
		contents,
		request.First && submission.Access.Owner,
	)
	if err != nil {
		return nil, err
	}

	if err := submission.SaveBeatmap(store, beatmap.Target.Id, contents); err != nil {
		return nil, err
	}
	return oszFile, nil
}

func handleOszUpload(request *PreOsz2UploadRequest, submission *bss.SubmissionContext, store storage.Storage) (*os.File, error) {
	if !strings.EqualFold(path.Ext(request.FileHeader.Filename), ".osz") {
		return nil, fmt.Errorf("%w: unsupported legacy upload %q", bss.ErrInvalidSubmissionRequest, request.FileHeader.Filename)
	}

	reader, err := zip.NewReader(request.File, request.Size())
	if err != nil {
		return nil, fmt.Errorf("%w: open osz: %v", bss.ErrInvalidSubmissionRequest, err)
	}
	if err := submission.SetFS(reader); err != nil {
		return nil, err
	}
	if err := submission.ValidateFS(); err != nil {
		return nil, fmt.Errorf("%w: validate osz: %v", bss.ErrInvalidSubmissionRequest, err)
	}

	oszFile, err := submission.BuildOsz(store)
	if err != nil {
		return nil, err
	}
	return oszFile, nil
}

func validateUploadTickets(secret string, userId int, request *PreOsz2UploadRequest) (int, error) {
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
