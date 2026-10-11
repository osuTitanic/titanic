package routes

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/osuTitanic/titanic/services/deck/internal/bss"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

type Osz2UploadRequest struct {
	Username       string
	Password       string
	SetId          int
	SubmissionType Osz2SubmissionType
	Osz2Checksum   string

	File       multipart.File
	FileHeader *multipart.FileHeader
	Form       *multipart.Form
}

func (request *Osz2UploadRequest) Size() int64 {
	if request.FileHeader == nil {
		return -1
	}
	return request.FileHeader.Size
}

func (request *Osz2UploadRequest) Close() {
	if request == nil {
		return
	}
	if request.File != nil {
		request.File.Close()
		request.File = nil
	}
	if request.Form != nil {
		request.Form.RemoveAll()
		request.Form = nil
	}
}

func NewOsz2UploadRequest(ctx *server.Context) (request *Osz2UploadRequest, err error) {
	err = ctx.EnsureMultipartForm(
		bss.MaxPackageFileSize,
		1<<20, // 1 MiB in-memory (larger file parts -> temporary storage)
	)
	if err != nil {
		if sizeError, ok := errors.AsType[*http.MaxBytesError](err); ok {
			return nil, fmt.Errorf("%w: multipart body exceeds %d bytes", bss.ErrOsz2TooLarge, sizeError.Limit)
		}
		return nil, fmt.Errorf("parse multipart form: %w", err)
	}
	request = &Osz2UploadRequest{Form: ctx.Request.MultipartForm}

	defer func(r *Osz2UploadRequest) {
		// Close only if there was an error
		// We want to keep the file stream open for the submission
		if err != nil {
			r.Close()
		}
	}(request)

	request.Username = strings.TrimSpace(ctx.FormValue("u"))
	request.Password = strings.TrimSpace(ctx.FormValue("h"))
	request.Osz2Checksum = strings.TrimSpace(ctx.FormValue("z"))

	request.SetId, err = ctx.FormValueInt("s")
	if err != nil {
		return nil, fmt.Errorf("parse beatmapset ID: %w", err)
	}
	if request.SetId <= 0 {
		return nil, fmt.Errorf("invalid beatmapset ID %d", request.SetId)
	}

	request.SubmissionType, err = ctx.FormValueEnum[Osz2SubmissionType]("t")
	if err != nil {
		return nil, fmt.Errorf("parse submission type: %w", err)
	}

	request.File, request.FileHeader, err = ctx.Request.FormFile("0")
	if err != nil {
		return nil, fmt.Errorf("open osz2 upload: %w", err)
	}
	if request.Size() > bss.MaxPackageFileSize {
		return nil, fmt.Errorf("%w: osz2 upload is %d bytes", bss.ErrOsz2TooLarge, request.Size())
	}
	if request.Size() <= 0 {
		return nil, fmt.Errorf("invalid osz2 upload size %d", request.Size())
	}
	return request, nil
}

// /web/osu-osz2-bmsubmit-upload.php -> Applies a full osz2 package or patch
// and updates the canonical beatmap resources and metadata.
func BeatmapSubmissionOsz2Upload(ctx *server.Context) {
	request, err := NewOsz2UploadRequest(ctx)
	if errors.Is(err, bss.ErrOsz2TooLarge) {
		renderBssError(ctx, err)
		return
	}
	if err != nil {
		renderBssError(ctx, fmt.Errorf("%w: %v", bss.ErrInvalidSubmissionRequest, err))
		return
	}
	defer request.Close()

	user, err := ctx.AuthenticateUser(
		request.Username,
		request.Password,
		false, // doesn't have to be connected to bancho
	)
	if errors.Is(err, server.ErrUserNotFound) || errors.Is(err, server.ErrInvalidPassword) {
		renderBssError(ctx, bss.ErrAuthenticationFailed)
		return
	}
	if err != nil {
		renderBssError(ctx, fmt.Errorf("authenticate user: %w", err))
		return
	}

	userPermissions, err := ctx.State.Permissions.Resolve(user.Id)
	if err != nil {
		renderBssError(ctx, fmt.Errorf("resolve permissions for user %d: %w", user.Id, err))
		return
	}

	submission := bss.NewSubmissionContext(ctx.Context())
	submission.User = user

	// Check if user is allowed to upload any beatmaps
	if err := submission.CheckSubmissionEligibility(
		ctx.State.Config.BeatmapSubmissionEnabled,
		userPermissions,
	); err != nil {
		renderBssError(ctx, err)
		return
	}

	// TODO: Submission implementation

	ctx.RenderText(http.StatusOK, "0")
}
