package routes

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Lekuruu/gosu/pkg/beatmaps"
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/permissions"
	"github.com/osuTitanic/titanic/internal/state"
	"github.com/osuTitanic/titanic/services/deck/internal/bss"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

// TODO: /web/osu-bmsubmit-getid2.php & /web/osu-bmsubmit-getid3.php
// TODO: /web/osu-bmsubmit-novideo.php
// TODO: /web/osu-bmsubmit-post[3-2].php

// GetIdAction describes where a submitted difficulty appears in the
// sequence of getid requests used to update a beatmapset.
type GetIdAction int

const (
	GetIdStandard GetIdAction = iota
	GetIdFirst
	GetIdLast
	GetIdSingle
)

func (action GetIdAction) Valid() bool {
	return action >= GetIdStandard && action <= GetIdSingle
}

func (action GetIdAction) StartsSubmission() bool {
	return action == GetIdFirst || action == GetIdSingle
}

func (action GetIdAction) FinishesSubmission() bool {
	return action == GetIdLast || action == GetIdSingle
}

// PreOsz2GetIdRequest contains a pre-osz2 getid request.
type PreOsz2GetIdRequest struct {
	Username      string
	Password      string
	SetId         int
	Action        GetIdAction
	HasVideo      bool
	HasStoryboard bool

	File       multipart.File
	FileHeader *multipart.FileHeader
}

func NewPreOsz2GetIdRequest(ctx *server.Context) (request *PreOsz2GetIdRequest, err error) {
	err = ctx.EnsureMultipartForm(
		bss.MaxPackageFileSize,
		0, // stream the uploaded .osu file to disk
	)
	if sizeError, ok := errors.AsType[*http.MaxBytesError](err); ok {
		return nil, fmt.Errorf("%w: multipart body exceeds %d bytes", bss.ErrOsz2TooLarge, sizeError.Limit)
	}
	if err != nil {
		return nil, fmt.Errorf("parse multipart form: %w", err)
	}

	request = &PreOsz2GetIdRequest{
		Username: ctx.FormValue("u"),
		Password: ctx.FormValue("p"),
		SetId:    -1,
	}

	request.Action, err = ctx.FormValueEnum[GetIdAction]("r")
	if err != nil {
		return nil, fmt.Errorf("parse getid action: %w", err)
	}

	setId, err := ctx.FormValueIntOptional("s")
	if err != nil {
		return nil, fmt.Errorf("parse beatmapset ID: %w", err)
	}
	if setId != nil {
		request.SetId = *setId
	}

	hasVideo, _ := ctx.FormValueIntOptional("v")
	hasStoryboard, _ := ctx.FormValueIntOptional("sb")

	request.HasVideo = hasVideo != nil && *hasVideo == 1
	request.HasStoryboard = hasStoryboard != nil && *hasStoryboard == 1

	request.File, request.FileHeader, err = ctx.FormFile("osu")
	if err != nil {
		return nil, fmt.Errorf("open beatmap upload: %w", err)
	}
	if request.Size() <= 0 {
		request.Close()
		return nil, fmt.Errorf("invalid beatmap upload size %d", request.Size())
	}
	if request.Size() > bss.MaxBeatmapFileSize {
		size := request.Size()
		request.Close()
		return nil, fmt.Errorf("%w: beatmap upload is %d bytes", bss.ErrOsz2TooLarge, size)
	}
	return request, nil
}

func (request *PreOsz2GetIdRequest) Size() int64 {
	if request.FileHeader == nil {
		return -1
	}
	return request.FileHeader.Size
}

func (request *PreOsz2GetIdRequest) Close() {
	if request == nil {
		return
	}
	if request.File == nil {
		return
	}
	request.File.Close()
	request.File = nil
}

func (request *PreOsz2GetIdRequest) ReadBeatmap() (beatmap bss.StagedSubmissionBeatmap, source *beatmaps.Beatmap, err error) {
	if request.File == nil {
		return beatmap, nil, errors.New("beatmap file is nil")
	}

	contents, err := io.ReadAll(io.LimitReader(request.File, bss.MaxBeatmapFileSize+1))
	if err != nil {
		return beatmap, nil, fmt.Errorf("read beatmap upload: %w", err)
	}
	if len(contents) > bss.MaxBeatmapFileSize {
		return beatmap, nil, fmt.Errorf(
			"%w: beatmap upload exceeds %d bytes",
			bss.ErrOsz2TooLarge,
			bss.MaxBeatmapFileSize,
		)
	}

	beatmap = bss.StagedSubmissionBeatmap{
		Id:       -1,
		Filename: request.FileHeader.Filename,
		Checksum: fmt.Sprintf("%x", md5.Sum(contents)),
		Contents: contents,
	}
	if err := beatmap.Validate(); err != nil {
		return beatmap, nil, fmt.Errorf("validate beatmap upload: %w", err)
	}
	source, err = beatmap.Parse()
	if err != nil {
		return beatmap, nil, fmt.Errorf("parse beatmap upload: %w", err)
	}
	return beatmap, source, nil
}

// PreOsz2GetIdResponse contains the response values
// returned for one submitted difficulty.
type PreOsz2GetIdResponse struct {
	SetId         int
	OszTicket     string
	BeatmapTicket string
	OszFilename   string

	Topic  *bss.TopicContents
	Status constants.BeatmapStatus
}

func (response PreOsz2GetIdResponse) FormatResponse() string {
	values := []string{
		"new",
		strconv.Itoa(response.SetId),
		response.OszTicket,
		response.BeatmapTicket,
		response.OszFilename,
	}
	if response.Topic == nil {
		return strings.Join(values, "\n")
	}

	values[0] = "old"
	values = append(values,
		strconv.Itoa(response.Topic.Id),
		strconv.Itoa(response.Status.Value()),
		response.Topic.Subject,
		response.Topic.Message,
	)
	return strings.Join(values, "\n")
}

// /web/osu-bmsubmit-getid[4-5].php -> Stages one .osu file & applies the
// complete submission when the client sent its last / single difficulty
func BeatmapSubmissionPreOsz2GetId(ctx *server.Context) {
	request, err := NewPreOsz2GetIdRequest(ctx)
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

	beatmap, source, err := request.ReadBeatmap()
	if errors.Is(err, bss.ErrOsz2TooLarge) {
		renderBssErrorText(ctx, err)
		return
	}
	if err != nil {
		renderBssErrorText(ctx, fmt.Errorf("%w: %v", bss.ErrInvalidSubmissionRequest, err))
		return
	}

	staged, err := ctx.StagedSubmissions.Get(
		ctx.Context(),
		user.Id,
	)
	if request.Action.StartsSubmission() {
		// This is the first beatmap upload of this submission
		// We will check if the beatmap is part of an existing
		// set, or if the user wants to create a new one
		staged, err = beginSubmission(
			ctx, request, submission,
			userPermissions, source,
		)
	}
	if err != nil {
		renderBssErrorText(ctx, err)
		return
	}
	if staged == nil {
		renderBssErrorText(ctx, errors.New("bss: staged submission not found"))
		return
	}

	if request.SetId > 0 && request.SetId != staged.BeatmapsetId {
		// User tried to update the wrong beatmapset / changed the
		// beatmapset mid-submission, which should not be possible
		renderBssErrorText(ctx, fmt.Errorf(
			"%w: request beatmapset %d does not match staged beatmapset %d",
			bss.ErrInvalidSubmissionRequest, request.SetId, staged.BeatmapsetId,
		))
		return
	}

	// Append the new beatmap to our staged submission
	// (saved in redis, later applied once we have all beatmaps)
	staged.Beatmaps = append(staged.Beatmaps, beatmap)
	staged.OszFilename = buildOszFilename(staged.BeatmapsetId, source.Artist, source.Title)

	if request.Action.FinishesSubmission() {
		// We have gathered all beatmaps!
		// Now we can start updating all the beatmap &
		// beatmapset metadata using the staged submission
		// Note that this won't update the osz files yet
		// (that will be done by the upload endpoint later)
		err := applySubmission(
			ctx, submission,
			userPermissions, staged,
		)
		if err != nil {
			renderBssErrorText(ctx, err)
			return
		}
	}

	if submission.Beatmapset == nil {
		submission.Beatmapset, err = ctx.State.Beatmapsets.ById(staged.BeatmapsetId)
		if err != nil {
			renderBssErrorText(ctx, fmt.Errorf("bss: fetch staged beatmapset %d: %w", staged.BeatmapsetId, err))
			return
		}
		if submission.Beatmapset == nil {
			renderBssErrorText(ctx, fmt.Errorf("%w: %d", bss.ErrBeatmapsetNotFound, staged.BeatmapsetId))
			return
		}
	}

	oszTicket, err := bss.NewSubmissionTicket(
		ctx.State.Config.FrontendSecretKey,
		user.Id,
		staged.BeatmapsetId,
		bss.SubmissionTicketOsz,
		staged.OszFilename,
	)
	if err != nil {
		renderBssErrorText(ctx, err)
		return
	}

	beatmapTicket, err := bss.NewSubmissionTicket(
		ctx.State.Config.FrontendSecretKey,
		user.Id,
		staged.BeatmapsetId,
		bss.SubmissionTicketBeatmap,
		beatmap.Filename,
	)
	if err != nil {
		renderBssErrorText(ctx, err)
		return
	}

	topic, err := bss.GetTopic(
		ctx.State.Repositories,
		submission.Beatmapset,
	)
	if err != nil {
		renderBssErrorText(ctx, err)
		return
	}
	if topic.Id == 0 {
		topic = nil
	}

	response := PreOsz2GetIdResponse{
		SetId:         staged.BeatmapsetId,
		OszTicket:     oszTicket,
		BeatmapTicket: beatmapTicket,
		OszFilename:   staged.OszFilename,
		Topic:         topic,
		Status:        submission.Beatmapset.Status,
	}

	if request.Action.FinishesSubmission() {
		if err := ctx.StagedSubmissions.Delete(ctx.Context(), user.Id); err != nil {
			ctx.Logger.Warn("Failed to delete staged beatmap submission", "user_id", user.Id, "error", err)
		}
	} else if err := ctx.StagedSubmissions.Save(ctx.Context(), staged); err != nil {
		renderBssErrorText(ctx, err)
		return
	}

	ctx.Logger.Info(
		"Handled beatmap submission",
		"user_id", user.Id,
		"set_id", staged.BeatmapsetId,
		"action", request.Action,
		"filename", beatmap.Filename,
	)
	ctx.RenderText(http.StatusOK, response.FormatResponse())
}

func beginSubmission(
	ctx *server.Context,
	request *PreOsz2GetIdRequest,
	submission *bss.SubmissionContext,
	userPermissions *permissions.Set,
	source *beatmaps.Beatmap,
) (*bss.StagedSubmission, error) {
	newBeatmapset := false

	err := ctx.State.DatabaseTransaction(func(transaction *state.Repositories) error {
		// We first want to remove any inactive maps
		removed, err := submission.RemoveInactiveBeatmapsets(
			transaction,
			ctx.State.Storage,
		)
		if err != nil {
			ctx.Logger.Warn("Failed to remove inactive beatmapsets", "error", err)
		}
		for _, beatmapset := range removed {
			ctx.Logger.Info("Removed inactive beatmapset", "id", beatmapset.Id, "name", beatmapset.Name())
			ctx.State.Resources.InvalidateCaches(ctx.Context(), beatmapset.Id)
		}

		existing, err := prepareBeatmapsetFromSource(transaction, submission, request.SetId, source)
		if err != nil {
			return err
		}

		if !existing {
			// We received a new submission!
			remainingSlots, err := submission.RemainingUploadSlots(transaction, userPermissions)
			if err != nil {
				return err
			}
			if remainingSlots <= 0 {
				return bss.ErrNoUploadSlots
			}
			// Create an empty / placeholder set for the new submission
			if err := submission.PrepareNewBeatmapset(transaction); err != nil {
				return err
			}
			newBeatmapset = true
		}

		if err := submission.ResolveAccess(transaction); err != nil {
			return err
		}
		if existing {
			// Check if we have the permission to update this set (as the owner or collaborator)
			return submission.CheckBeatmapsetEligibility()
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	staged := &bss.StagedSubmission{
		UserId:        submission.User.Id,
		BeatmapsetId:  submission.Beatmapset.Id,
		NewBeatmapset: newBeatmapset,
		HasVideo:      request.HasVideo,
		HasStoryboard: request.HasStoryboard,
		Beatmaps:      make([]bss.StagedSubmissionBeatmap, 0),
	}
	return staged, nil
}

func applySubmission(
	ctx *server.Context,
	submission *bss.SubmissionContext,
	userPermissions *permissions.Set,
	staged *bss.StagedSubmission,
) error {
	return ctx.State.DatabaseTransaction(func(transaction *state.Repositories) error {
		if err := submission.PrepareBeatmapset(transaction, staged.BeatmapsetId); err != nil {
			return err
		}
		if err := submission.ResolveAccess(transaction); err != nil {
			return err
		}
		// Check if we have the permission to update this set (as the owner or collaborator)
		if err := submission.CheckBeatmapsetEligibility(); err != nil {
			return err
		}

		// Assign all staged beatmaps to the submission & create set metadata
		// This will also check what beatmaps got removed & which ones are newly added
		if err := submission.PrepareStagedBeatmaps(
			transaction,
			staged,
			ctx.State.PPv1,
			ctx.State.PPv2,
		); err != nil {
			return err
		}

		// This will later be further validated in the upload endpoint
		if submission.MaxBeatmapLength() <= 1 {
			return bss.ErrBeatmapTooShort
		}

		if submission.Access.Owner {
			// Remove beatmaps from the database that are not in this submission
			if err := submission.RemoveBeatmaps(transaction, submission.RemovedBeatmaps()); err != nil {
				return err
			}
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

		// Determine & persist new beatmapset & beatmap metadata to the database
		// Will be committed once we exit the transaction without an error

		// Beatmapset metadata (e.g. title, artist, creator)
		// should only be updated by the set owner
		if submission.Access.Owner {
			if err := bss.ApplyBeatmapsetMetadata(transaction, submission.Beatmapset, submission.Metadata); err != nil {
				return err
			}
		}
		// Finally, update each individual difficulty w/ their derived metadata
		for _, beatmap := range submission.Beatmaps {
			if beatmap.Source == nil {
				continue
			}
			if err := bss.ApplyBeatmapMetadata(transaction, beatmap.Target); err != nil {
				return err
			}
		}

		columns := []string{"last_updated"}
		submission.Beatmapset.LastUpdate = time.Now()

		if submission.Access.Owner {
			// "PrepareStagedBeatmaps" may have updated the video / sb flags
			columns = append(columns, "has_video")
			columns = append(columns, "has_storyboard")
		}

		_, err := transaction.Beatmapsets.Update(submission.Beatmapset, columns...)
		if err != nil {
			return fmt.Errorf("bss: update staged beatmapset %d: %w", submission.Beatmapset.Id, err)
		}

		// "PopBubble" essentially means removing all BAT nominations
		return bss.PopBubble(transaction, submission.Beatmapset)
	})
}

func prepareBeatmapsetFromSource(
	transaction *state.Repositories,
	submission *bss.SubmissionContext,
	requestedSetId int,
	source *beatmaps.Beatmap,
) (bool, error) {
	if requestedSetId > 0 {
		err := submission.PrepareBeatmapset(transaction, requestedSetId)
		if err == nil {
			// We got a valid beatmapset, either existing or new
			return true, nil
		}
		if !errors.Is(err, bss.ErrBeatmapsetNotFound) {
			// No beatmapset was found with that ID
			return false, err
		}
		// welp, that didn't work
		// let's try to find it by artist, title & creator
	}

	beatmapset, err := transaction.Beatmapsets.ByArtistTitleAndCreator(
		source.Artist,
		source.Title,
		submission.User.Id,
	)
	if err != nil {
		return false, fmt.Errorf("bss: find beatmapset by artist, title and owner: %w", err)
	}
	if beatmapset == nil {
		return false, nil
	}
	if err := submission.PrepareBeatmapset(transaction, beatmapset.Id); err != nil {
		return false, err
	}
	return true, nil
}

func buildOszFilename(setId int, artist, title string) string {
	filename := fmt.Sprintf(
		"%d %s - %s.osz",
		setId, artist, title,
	)
	filter := func(character rune) rune {
		if strings.ContainsRune(`\/:*?"<>|`, character) {
			return -1
		}
		return character
	}
	filename = strings.Map(filter, filename)
	return strings.ReplaceAll(filename, "#", "%23")
}

func renderBssErrorText(ctx *server.Context, err error) {
	// Ensure everything was read from the request body
	// Otherwise the client will not display any errors!
	_, drainErr := io.Copy(io.Discard, ctx.Request.Body)
	if drainErr != nil {
		ctx.Logger.Warn("Failed to drain request body", "err", drainErr)
	}

	message, expected := bss.FormatPreOsz2Error(err)
	if expected {
		ctx.Logger.Warn("Rejected pre-osz2 submission", "error", err)
	} else {
		ctx.Logger.Error("Failed to process pre-osz2 submission", "error", err)
	}
	ctx.RenderText(http.StatusOK, message)
}

// TODO: Check if we can consolidate renderBssError & renderBssErrorText
