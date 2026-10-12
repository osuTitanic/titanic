package bss

import (
	"fmt"
	"strings"
	"time"

	"github.com/Lekuruu/osz2-go/pkg/osz2"
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

// RequireBeatmapsetOwner ensures that the authenticated user owns the beatmapset.
func (submission *SubmissionContext) RequireBeatmapsetOwner() error {
	if submission.User == nil {
		return ErrUserNotAuthenticated
	}
	if submission.Beatmapset == nil {
		return ErrBeatmapsetNotPrepared
	}
	if submission.Beatmapset.CreatorId == nil || *submission.Beatmapset.CreatorId != submission.User.Id {
		return ErrBeatmapsetAccessNotAllowed
	}
	return nil
}

// ValidateOwnership ensures that creator names belong to the beatmapset owner or one of the collaborators.
func (submission *SubmissionContext) ValidateOwnership(allowedUsernames []string) error {
	allowed := make(map[string]bool, len(allowedUsernames))

	for _, username := range allowedUsernames {
		username = strings.ToLower(strings.TrimSpace(username))
		if username != "" {
			allowed[username] = true
		}
	}

	isAllowed := func(username string) bool {
		_, ok := allowed[strings.ToLower(strings.TrimSpace(username))]
		return ok
	}

	// Check osz2 metadata
	if creator := submission.Metadata[osz2.Creator]; creator != "" && !isAllowed(creator) {
		return ErrBeatmapOwnershipNotAllowed
	}
	// Check beatmap files
	for _, beatmap := range submission.Beatmaps {
		if beatmap.Source != nil && !isAllowed(beatmap.Source.Creator) {
			return ErrBeatmapOwnershipNotAllowed
		}
	}
	return nil
}

// AllowedUsernames returns all usernames that are valid for submission.
// This includes the set owner, set collaborators and all of their previous usernames.
func (submission *SubmissionContext) AllowedUsernames(transaction *state.Repositories) ([]string, error) {
	// Add current user that is submitting
	userIds := map[int]bool{submission.User.Id: true}
	usernames := []string{submission.User.Name}

	// Add set owner
	if submission.Beatmapset.Creator != nil {
		usernames = append(usernames, *submission.Beatmapset.Creator)
	}
	if submission.Beatmapset.CreatorId != nil {
		userIds[*submission.Beatmapset.CreatorId] = true
	}

	// Collect all beatmap IDs to fetch collaborators
	beatmapIds := make([]int, len(submission.Beatmaps))
	for i, beatmap := range submission.Beatmaps {
		beatmapIds[i] = beatmap.Target.Id
	}

	collaborations, err := transaction.Collaborations.FetchByBeatmaps(beatmapIds, "User")
	if err != nil {
		return nil, fmt.Errorf("bss: fetch collaborators for set %d: %w", submission.Beatmapset.Id, err)
	}

	// Add all collaborators for this set
	for _, collaboration := range collaborations {
		userIds[collaboration.UserId] = true
		usernames = append(usernames, collaboration.User.Name)
	}

	// For each user that we collected, fetch their previous names
	// from the name history & add it to the allowed list
	for userId := range userIds {
		user, err := transaction.Users.ById(userId)
		if err != nil {
			return nil, fmt.Errorf("bss: fetch submission participant %d: %w", userId, err)
		}
		if user != nil {
			usernames = append(usernames, user.Name)
		}

		names, err := transaction.Names.ManyByUserId(userId)
		if err != nil {
			return nil, fmt.Errorf("bss: fetch name history for submission participant %d: %w", userId, err)
		}
		for _, name := range names {
			usernames = append(usernames, name.Name)
		}
	}
	return usernames, nil
}
