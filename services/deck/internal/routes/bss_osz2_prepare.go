package routes

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
	"github.com/osuTitanic/titanic/services/deck/internal/bss"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

const (
	osz2SubmissionFull  = 1
	osz2SubmissionPatch = 2
)

type Osz2GetIdRequest struct {
	SetId        int
	BeatmapIds   []int
	Osz2Checksum string
}

type Osz2GetIdResult struct {
	SetId          int
	BeatmapIds     []int
	SubmissionType int
	RemainingSlots *int
	Bubbled        bool
	NewBeatmapset  bool
}

func (result Osz2GetIdResult) FormatResponse() string {
	beatmapIds := make([]string, len(result.BeatmapIds))
	for i, beatmapId := range result.BeatmapIds {
		beatmapIds[i] = strconv.Itoa(beatmapId)
	}

	remainingSlots := ""
	if result.RemainingSlots != nil {
		remainingSlots = strconv.Itoa(*result.RemainingSlots)
	}

	bubbled := "0"
	if result.Bubbled {
		bubbled = "1"
	}

	return strings.Join([]string{
		"0", // status / error code (0 = success)
		strconv.Itoa(result.SetId),
		strings.Join(beatmapIds, ","),
		strconv.Itoa(result.SubmissionType),
		remainingSlots,
		bubbled,
	}, "\n")
}

func NewOsz2GetIdRequest(ctx *server.Context) (request Osz2GetIdRequest, err error) {
	osz2Checksum := strings.TrimSpace(ctx.QueryValue("z")) // may be empty

	setId, err := ctx.QueryValueInt("s") // -1 if new submission
	if err != nil {
		return request, fmt.Errorf("parse beatmapset ID: %w", err)
	}

	beatmapIdsRaw := strings.Split(ctx.QueryValue("b"), ",") // items are also -1 if the map is new
	beatmapIds := make([]int, len(beatmapIdsRaw))

	for i, rawBeatmapId := range beatmapIdsRaw {
		beatmapId, err := strconv.Atoi(strings.TrimSpace(rawBeatmapId))
		if err != nil {
			return request, fmt.Errorf("parse beatmap ID %q: %w", rawBeatmapId, err)
		}
		beatmapIds[i] = beatmapId
	}

	return Osz2GetIdRequest{
		SetId:        setId,
		BeatmapIds:   beatmapIds,
		Osz2Checksum: osz2Checksum,
	}, nil
}

// /web/osu-osz2-bmsubmit-getid.php -> Prepares an osz2 submission
// and generates its server-side beatmapset & beatmap IDs.
func BeatmapSubmissionOsz2GetId(ctx *server.Context) {
	request, err := NewOsz2GetIdRequest(ctx)
	if err != nil {
		renderBssError(ctx, fmt.Errorf("%w: %v", bss.ErrInvalidSubmissionRequest, err))
		return
	}

	user, err := ctx.AuthenticateUserFromQuery("u", "h", false)
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

	// Check if the user is allowed to upload a beatmapset
	if err := submission.CheckSubmissionEligibility(
		ctx.State.Config.BeatmapSubmissionEnabled,
		userPermissions,
	); err != nil {
		renderBssError(ctx, err)
		return
	}

	// Actual processing of the request happens now
	var result Osz2GetIdResult

	err = ctx.State.DatabaseTransaction(func(transaction *state.Repositories) error {
		// We first want to remove any inactive maps
		if _, err := submission.RemoveInactiveBeatmapsets(transaction); err != nil {
			return err
		}

		existingBeatmapset, err := prepareBeatmapset(transaction, submission, request.SetId)
		if err != nil {
			return err
		}

		if !existingBeatmapset {
			// We received a new submission!
			submission.Beatmapset = nil
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
			result.RemainingSlots = &remainingSlots
		}

		if err := submission.ResolveAccess(transaction); err != nil {
			return err
		}
		if existingBeatmapset {
			// Check if we have the permission to update this set (as the owner or collaborator)
			if err := submission.CheckBeatmapsetEligibility(); err != nil {
				return err
			}
		}

		// Check what beatmaps got removed & which ones are newly added
		if err := submission.PrepareBeatmaps(transaction, request.BeatmapIds); err != nil {
			return err
		}
		if err := submission.RemoveBeatmaps(transaction, submission.RemovedBeatmaps()); err != nil {
			return err
		}

		bubbled, err := isBubbledBeatmapset(transaction, submission.Beatmapset)
		if err != nil {
			return err
		}

		result.SetId = submission.Beatmapset.Id
		result.BeatmapIds = preparedBeatmapIds(submission)
		result.Bubbled = bubbled
		result.NewBeatmapset = !existingBeatmapset
		return nil
	})
	if err != nil {
		renderBssError(ctx, err)
		return
	}

	result.SubmissionType = resolveOsz2SubmissionType(
		ctx,
		result.SetId,
		result.NewBeatmapset,
		request.Osz2Checksum,
	)

	ctx.Logger.Info(
		"Prepared beatmap submission, ready for upload",
		"user_id", user.Id,
		"set_id", result.SetId,
		"beatmap_ids", result.BeatmapIds,
		"submission_type", result.SubmissionType,
	)
	ctx.RenderText(http.StatusOK, result.FormatResponse())
}

func prepareBeatmapset(transaction *state.Repositories, submission *bss.SubmissionContext, setId int) (bool, error) {
	if setId <= 0 {
		return false, nil
	}

	err := submission.PrepareBeatmapset(transaction, setId)
	if err == nil {
		// We got a valid beatmapset, either existing or new
		return submission.Beatmapset.Status != constants.BeatmapStatusInactive, nil
	}
	if errors.Is(err, bss.ErrBeatmapsetNotFound) {
		// No beatmapset was found with that ID
		return false, nil
	}
	// shit happened
	return false, err
}

func preparedBeatmapIds(submission *bss.SubmissionContext) []int {
	// Return the new beatmap IDs in the correct order
	// The client expects them that way, otherwise the IDs would be shuffled around
	beatmapIds := make([]int, len(submission.Beatmaps))
	for i, beatmap := range submission.Beatmaps {
		beatmapIds[i] = beatmap.Target.Id
	}
	return beatmapIds
}

func isBubbledBeatmapset(transaction *state.Repositories, beatmapset *schemas.Beatmapset) (bool, error) {
	// Bubbled means the set got at least one nomination
	// We could check the topic icon here, but I think this is cleaner
	isBubbled, err := transaction.Nominations.ExistsForSet(beatmapset.Id)
	if err != nil {
		return false, fmt.Errorf("fetch nomination status: %w", err)
	}
	return isBubbled, nil
}

func resolveOsz2SubmissionType(ctx *server.Context, setId int, newBeatmapset bool, clientChecksum string) int {
	if newBeatmapset || clientChecksum == "" {
		return osz2SubmissionFull
	}
	if ctx.State.StorageOsz2 == nil {
		return osz2SubmissionFull
	}

	// We use this to determine if the client can upload a patch for the
	// existing osz2 or if the client has to upload the full osz2 instead

	stream, err := ctx.State.StorageOsz2.ReadStream(ctx.Context(), strconv.Itoa(setId), "osz2")
	if err != nil {
		ctx.Logger.Debug("Unable to read server-side osz2", "set_id", setId, "error", err)
		return osz2SubmissionFull
	}
	defer stream.Close()

	hash := md5.New()
	_, err = io.Copy(hash, stream)

	if err != nil {
		ctx.Logger.Warn("Failed to hash server-side osz2", "set_id", setId, "error", err)
		return osz2SubmissionFull
	}
	serverChecksum := hex.EncodeToString(hash.Sum(nil))

	if !strings.EqualFold(clientChecksum, serverChecksum) {
		// Server-side osz2 doesn't match client-side osz2
		return osz2SubmissionFull
	}

	// We can safely receive a patch from the client!
	return osz2SubmissionPatch
}

func renderBssError(ctx *server.Context, err error) {
	response, expected := bss.FormatOsz2ErrorResponse(err)
	if expected {
		ctx.Logger.Warn("Rejected osz2 submission", "error", err)
	} else {
		ctx.Logger.Error("Failed to process osz2 submission", "error", err)
	}
	ctx.RenderText(http.StatusOK, response)
}
