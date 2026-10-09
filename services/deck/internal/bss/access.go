package bss

import (
	"fmt"
	"time"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/permissions"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
)

// CheckSubmissionEligibility ensures that the user is eligible to submit beatmaps.
func (submission *SubmissionContext) CheckSubmissionEligibility(submissionEnabled bool, userPermissions *permissions.Set) error {
	if !submissionEnabled {
		return ErrSubmissionDisabled
	}
	if submission.User == nil {
		return ErrUserNotAuthenticated
	}
	if !submission.User.Activated {
		return ErrUserNotActivated
	}
	if submission.User.Restricted {
		return ErrUserRestricted
	}
	if submission.User.SilenceEnd != nil && submission.User.SilenceEnd.After(time.Now()) {
		return ErrUserSilenced
	}
	if submission.User.IsBot {
		return ErrUserIsBot
	}
	if !userPermissions.Has("beatmaps.upload") {
		return ErrUploadPermissionDenied
	}
	return nil
}

// CheckBeatmapsetEligibility ensures that the beatmapset can be updated.
func (submission *SubmissionContext) CheckBeatmapsetEligibility() error {
	if submission.Beatmapset == nil {
		return ErrBeatmapsetNotPrepared
	}
	if submission.Beatmapset.Server != constants.BeatmapServerTitanic {
		return ErrBeatmapsetWrongServer
	}
	if submission.Beatmapset.Status >= constants.BeatmapStatusRanked {
		return ErrBeatmapsetRanked
	}
	if submission.Beatmapset.Status == constants.BeatmapStatusGraveyard {
		return ErrBeatmapsetGraveyarded
	}
	return nil
}

// RemainingUploadSlots returns how many new beatmapsets the user can upload.
func (submission *SubmissionContext) RemainingUploadSlots(transaction *state.Repositories, userPermissions *permissions.Set) (int, error) {
	if userPermissions.IsBat() {
		return 99, nil
	}

	unranked, err := transaction.Beatmapsets.CountUnrankedByCreator(submission.User.Id)
	if err != nil {
		return 0, fmt.Errorf("bss: count unranked beatmapsets for user %d: %w", submission.User.Id, err)
	}
	ranked, err := transaction.Beatmapsets.CountRankedByCreator(submission.User.Id)
	if err != nil {
		return 0, fmt.Errorf("bss: count ranked beatmapsets for user %d: %w", submission.User.Id, err)
	}

	baseLimit := 4
	rankedBonusLimit := 4
	if userPermissions.Has("beatmaps.upload.extended_limit") {
		baseLimit = 8
		rankedBonusLimit = 12
	}
	return baseLimit - unranked + min(ranked, rankedBonusLimit), nil
}

/*
 * Unlike the official beatmap submission system, Titanic! has its own collaboration
 * system, allowing updates on individual difficulties, if the user is a collaborator.
 *
 * Collaborators may only update their assigned beatmaps. Only with a separate permission
 * they are able to update beatmapset resources (e.g. sb files).
 */

// SubmissionAccess describes which parts of
// a beatmapset the user may update.
type SubmissionAccess struct {
	Owner              bool
	BeatmapIds         map[int]bool
	CanUpdateResources bool
}

func (access *SubmissionAccess) CanUpdateBeatmap(beatmapId int) bool {
	if access == nil {
		return false
	}
	if access.Owner {
		return true
	}
	_, ok := access.BeatmapIds[beatmapId]
	return ok
}

// ResolveAccess resolves beatmap & resource access for the authenticated user.
func (submission *SubmissionContext) ResolveAccess(transaction *state.Repositories) error {
	owner := submission.Beatmapset.CreatorId != nil &&
		*submission.Beatmapset.CreatorId == submission.User.Id
	access := &SubmissionAccess{
		Owner:              owner,
		CanUpdateResources: owner,
		BeatmapIds:         make(map[int]bool),
	}

	if owner {
		// thoust hadst allmighty power
		submission.Access = access
		return nil
	}

	beatmapIds := make([]int, 0, len(submission.Beatmapset.Beatmaps))
	for _, beatmap := range submission.Beatmapset.Beatmaps {
		beatmapIds = append(beatmapIds, beatmap.Id)
	}

	collaborations, err := transaction.Collaborations.FetchByUserAndBeatmapsWithLock(
		submission.User.Id,
		beatmapIds,
	)
	if err != nil {
		return fmt.Errorf("bss: resolve collaboration access for set %d: %w", submission.Beatmapset.Id, err)
	}
	if len(collaborations) == 0 {
		return ErrBeatmapsetAccessNotAllowed
	}

	for _, collaboration := range collaborations {
		access.BeatmapIds[collaboration.BeatmapId] = true
		access.CanUpdateResources = access.CanUpdateResources || collaboration.AllowResourceUpdates
	}
	submission.Access = access
	return nil
}

// CheckBeatmapAccess ensures that the authenticated user may update a difficulty.
func (submission *SubmissionContext) CheckBeatmapAccess(beatmap *schemas.Beatmap) error {
	if beatmap.SetId != submission.Beatmapset.Id {
		return ErrBeatmapAccessNotAllowed
	}
	if !submission.Access.CanUpdateBeatmap(beatmap.Id) {
		return ErrBeatmapAccessNotAllowed
	}
	return nil
}

// RequireResourceAccess ensures that the user is allowed to update beatmapset resources.
func (submission *SubmissionContext) RequireResourceAccess() error {
	if submission.Access == nil || !submission.Access.CanUpdateResources {
		return ErrResourceAccessNotAllowed
	}
	return nil
}
