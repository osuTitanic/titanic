package routes

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Lekuruu/osz2-go/pkg/osz2"
	"github.com/osuTitanic/titanic/internal/state"
	"github.com/osuTitanic/titanic/internal/storage"
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

func NewOsz2UploadRequest(ctx *server.Context) (request *Osz2UploadRequest, err error) {
	err = ctx.EnsureMultipartForm(
		bss.MaxPackageFileSize,
		0, // <- stream multipart data directly to a file instead of memory
	)
	if sizeError, ok := errors.AsType[*http.MaxBytesError](err); ok {
		return nil, fmt.Errorf("%w: multipart body exceeds %d bytes", bss.ErrOsz2TooLarge, sizeError.Limit)
	}
	if err != nil {
		return nil, fmt.Errorf("parse multipart form: %w", err)
	}
	request = &Osz2UploadRequest{Form: ctx.Request.MultipartForm}

	request.Username = strings.TrimSpace(ctx.FormValue("u"))
	request.Password = strings.TrimSpace(ctx.FormValue("h"))
	request.Osz2Checksum = strings.TrimSpace(ctx.FormValue("z"))

	if request.Username == "" {
		return nil, fmt.Errorf("empty username")
	}
	if request.Password == "" {
		return nil, fmt.Errorf("empty password")
	}

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

	request.File, request.FileHeader, err = ctx.FormFile("0", "osz2") // both "0" and "osz2" are filenames sent by the client
	if err != nil {
		request.Close()
		return nil, fmt.Errorf("open osz2 upload: %w", err)
	}
	if request.Size() > bss.MaxPackageFileSize {
		request.Close()
		return nil, fmt.Errorf("%w: osz2 upload is %d bytes", bss.ErrOsz2TooLarge, request.Size())
	}
	if request.Size() <= 0 {
		request.Close()
		return nil, fmt.Errorf("invalid osz2 upload size %d", request.Size())
	}

	return request, nil
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

	err = ctx.State.DatabaseTransaction(func(transaction *state.Repositories) error {
		err := submission.PrepareBeatmapset(transaction, request.SetId)
		if errors.Is(err, bss.ErrBeatmapsetNotFound) {
			return fmt.Errorf("%w: %d", bss.ErrReservedSetNotFound, request.SetId)
		}
		if err != nil {
			return err
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

		var osz2Source io.ReaderAt = request.File
		var osz2Size = request.Size()

		if request.SubmissionType == osz2SubmissionPatch {
			// We have received a patch file for the existing server-side osz2
			// This is used to speed up submission time & save bandwidth

			// We will use the osz2 base & the patch to create the new osz2
			// that contains the actual submission (streamed into a temporary file)
			// If the patch is successful, we can use it as our new base osz2 later on
			patchedFile, patchedSize, err := resolvePatchedOsz2(
				submission.Context,
				ctx.State.StorageOsz2,
				request.SetId,
				request.File,
				request.Size(),
			)
			if err != nil {
				return err
			}
			defer func() {
				patchedFile.Close()
				os.Remove(patchedFile.Name())
			}()

			osz2Source = patchedFile
			osz2Size = patchedSize
		}

		osz2Reader, err := osz2.NewReader(
			osz2Source,
			osz2Size,
			osz2.KeyTypeOsz2,
		)
		if err != nil {
			return fmt.Errorf("%w: open package for set %d: %v", bss.ErrOsz2Invalid, request.SetId, err)
		}

		if err := osz2Reader.Verify(); err != nil {
			return fmt.Errorf("%w: verify package for set %d: %v", bss.ErrOsz2Invalid, request.SetId, err)
		}
		if err := submission.SetFS(osz2Reader); err != nil {
			return err
		}
		if err := submission.ValidateFS(); err != nil {
			return fmt.Errorf("%w: validate package filesystem: %v", bss.ErrOsz2Invalid, err)
		}
		submission.Metadata = osz2Reader.Metadata()

		// Check if the osz2 lines up with the submission request data

		if err := validateOsz2BeatmapsetId(submission.Metadata, submission.Beatmapset.Id); err != nil {
			return err
		}
		if err := validateOsz2BeatmapIds(osz2Reader, submission.Beatmaps); err != nil {
			return err
		}

		if err := populateOsz2Beatmaps(ctx, submission, osz2Reader); err != nil {
			return err
		}

		if !userPermissions.IsBat() {
			allowedUsernames, err := submission.AllowedUsernames(transaction)
			if err != nil {
				return err
			}
			if err := submission.ValidateOwnership(allowedUsernames); err != nil {
				return err
			}
		}

		if err := submission.ValidateBeatmapUniqueness(transaction); err != nil {
			return err
		}

		// Build the new osz file, accounting for collaboration access
		oszFile, err := submission.BuildOsz(ctx.State.Storage)
		if err != nil {
			return err
		}
		// Cleanup the temporary osz after the transaction
		defer func() {
			oszFile.Close()
			os.Remove(oszFile.Name())
		}()

		oszInfo, err := oszFile.Stat()
		if err != nil {
			return fmt.Errorf("inspect osz: %w", err)
		}

		oszSize := oszInfo.Size()
		maxBeatmapLength := submission.MaxBeatmapLength()

		// The allowed osz size changes based on the beatmap length
		// The longer the beatmap, the higher the allowed package size
		if err := bss.ValidatePackageSize(oszSize, maxBeatmapLength); err != nil {
			return err
		}

		// Determine & persist new beatmapset & beatmap metadata to the database
		// Will be committed once we exit the transaction without an error
		if err := bss.ApplyBeatmapsetMetadata(transaction, submission.Beatmapset, submission.Metadata); err != nil {
			return err
		}
		if err := bss.ApplySubmissionMetadata(transaction, submission, oszFile); err != nil {
			return err
		}
		for _, beatmap := range submission.Beatmaps {
			if err := bss.ApplyBeatmapMetadata(transaction, beatmap.Target); err != nil {
				return err
			}
		}

		// Update new .osz file
		err = ctx.State.Storage.SaveStream(
			ctx.Context(),
			strconv.Itoa(submission.Beatmapset.Id),
			"osz",
			oszFile,
		)
		if err != nil {
			return err
		}

		// Upload all new .osu files
		for _, beatmap := range submission.Beatmaps {
			err = ctx.State.Storage.Save(
				ctx.Context(),
				strconv.Itoa(beatmap.Target.Id),
				"beatmaps",
				beatmap.Contents,
			)
		}

		// Upload new osz2 base that we can use for patching
		err = ctx.State.StorageOsz2.SaveStream(
			ctx.Context(),
			strconv.Itoa(submission.Beatmapset.Id),
			"osz2",
			io.NewSectionReader(osz2Source, 0, osz2Size),
		)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		renderBssError(ctx, err)
		return
	}

	// Post-commit stuff:
	// TODO: Generate thumbnail & audio preview
	// TODO: Broadcast events
	// TODO: Log success

	ctx.RenderText(http.StatusOK, "0")
}

func resolvePatchedOsz2(
	ctx context.Context,
	store storage.Storage,
	setId int,
	patch io.ReaderAt,
	patchSize int64,
) (output *os.File, size int64, err error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}
	if store == nil {
		return nil, 0, bss.ErrOsz2NotFound
	}

	source, sourceSize, err := store.ReadStreamAt(ctx, strconv.Itoa(setId), "osz2")
	if err != nil {
		return nil, 0, fmt.Errorf("%w: read cached package for set %d: %v", bss.ErrOsz2NotFound, setId, err)
	}
	defer source.Close()

	output, err = os.CreateTemp("", "bss-osz2-*")
	if err != nil {
		return nil, 0, fmt.Errorf("create temporary osz2: %w", err)
	}
	defer func() {
		// Only close & remove the file if an error occured
		// We want to keep the file open otherwise
		if err != nil {
			err = errors.Join(err, output.Close(), os.Remove(output.Name()))
			output = nil
		}
	}()

	size, err = osz2.ApplyPatch(
		source,
		sourceSize,
		patch,
		patchSize,
		bss.MaxPackageFileSize,
		output,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: apply patch for set %d: %v", bss.ErrOsz2Invalid, setId, err)
	}
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}
	return output, size, nil
}

func validateOsz2BeatmapsetId(metadata osz2.Metadata, expectedId int) error {
	value, ok := metadata[osz2.BeatmapSetID]
	if !ok {
		return fmt.Errorf("%w: missing beatmapset ID metadata", bss.ErrOsz2Invalid)
	}

	beatmapsetId, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("%w: invalid beatmapset ID metadata %q: %v", bss.ErrOsz2Invalid, value, err)
	}
	if beatmapsetId != expectedId {
		return fmt.Errorf(
			"%w: package beatmapset ID %d does not match expected ID %d",
			bss.ErrOsz2Invalid, beatmapsetId, expectedId,
		)
	}
	return nil
}

func validateOsz2BeatmapIds(reader *osz2.Reader, beatmaps []*bss.PreparedBeatmap) error {
	expectedIds := make(map[int32]bool, len(beatmaps))

	for _, beatmap := range beatmaps {
		beatmapId := int32(beatmap.Target.Id) // should hopefully never overflow...
		expectedIds[beatmapId] = true
	}

	// We want to check if the expected IDs line up with
	// whats inside the osz2 and that nothing is missing

	err := fs.WalkDir(reader, ".", func(filename string, entry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk osz2 file %q: %w", filename, err)
		}
		if entry.IsDir() {
			return nil
		}

		packageEntry, err := reader.Entry(filename)
		if err != nil {
			return fmt.Errorf("inspect osz2 file %q: %w", filename, err)
		}
		if !packageEntry.IsBeatmap() {
			return nil
		}

		beatmapId, assigned := packageEntry.BeatmapID()
		if !assigned {
			return fmt.Errorf("%w: beatmap %q has no assigned ID", bss.ErrOsz2Invalid, filename)
		}
		if _, expected := expectedIds[beatmapId]; !expected {
			return fmt.Errorf("%w: beatmap %q has unexpected ID %d", bss.ErrOsz2Invalid, filename, beatmapId)
		}
		delete(expectedIds, beatmapId)
		return nil
	})
	if err != nil {
		return err
	}

	// Ideally, no beatmaps should remain in expectedIds now
	// If there are, the osz2 is missing some beatmaps

	for _, beatmap := range beatmaps {
		beatmapId := int32(beatmap.Target.Id)
		_, missing := expectedIds[beatmapId]
		if missing {
			return fmt.Errorf("%w: package is missing beatmap ID %d", bss.ErrOsz2Invalid, beatmapId)
		}
	}
	return nil
}

func populateOsz2Beatmaps(ctx *server.Context, submission *bss.SubmissionContext, reader *osz2.Reader) error {
	if submission.Access == nil {
		return errors.New("assign osz2 beatmaps: submission access not resolved")
	}

	for _, beatmap := range submission.Beatmaps {
		beatmapId := beatmap.Target.Id
		canUpdate := submission.Access.CanUpdateBeatmap(beatmapId)
		if !canUpdate {
			continue
		}

		entry, exists := reader.EntryForBeatmap(int32(beatmapId))
		if !exists {
			return fmt.Errorf("%w: package is missing beatmap ID %d", bss.ErrOsz2Invalid, beatmapId)
		}

		_, _, err := submission.AssignBeatmapFile(
			beatmapId,
			entry.Path(),
		)
		if err != nil {
			return fmt.Errorf("%w: assign beatmap ID %d: %v", bss.ErrOsz2Invalid, beatmapId, err)
		}

		if err := beatmap.CalculateDifficulty(ctx.State.PPv1, ctx.State.PPv2); err != nil {
			return fmt.Errorf("calculate difficulty for beatmap %d: %w", beatmap.Target.Id, err)
		}
	}
	return nil
}
