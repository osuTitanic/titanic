package routes

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/state"
	"github.com/osuTitanic/titanic/services/deck/internal/bss"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

type Osz2GetTopicRequest struct {
	Username string
	Password string
	SetId    int
}

func NewOsz2GetTopicRequest(ctx *server.Context) (request Osz2GetTopicRequest, err error) {
	request.Username = ctx.QueryValue("u")
	request.Password = ctx.QueryValue("h")
	request.SetId, err = ctx.QueryValueInt("s")
	if err != nil {
		return request, fmt.Errorf("parse beatmapset ID: %w", err)
	}
	if request.SetId <= 0 {
		return request, fmt.Errorf("invalid beatmapset ID %d", request.SetId)
	}
	return request, nil
}

type Osz2PostTopicRequest struct {
	Username  string
	Password  string
	SetId     int
	Subject   string
	Message   string
	Complete  bool
	Subscribe bool
}

func NewOsz2PostTopicRequest(ctx *server.Context) (request Osz2PostTopicRequest, err error) {
	request.Username = ctx.FormValue("u")
	request.Password = ctx.FormValue("p")
	request.Subject = ctx.FormValue("subject")
	request.Message = ctx.FormValue("message")

	request.SetId, err = ctx.FormValueInt("b")
	if err != nil {
		return request, fmt.Errorf("parse beatmapset ID: %w", err)
	}
	if request.SetId <= 0 {
		return request, fmt.Errorf("invalid beatmapset ID %d", request.SetId)
	}

	complete, _ := ctx.FormValueInt("complete")
	subscribe, _ := ctx.FormValueInt("notify")
	request.Subscribe = subscribe == 1
	request.Complete = complete == 1

	return request, nil
}

func (request Osz2PostTopicRequest) TopicUpdate() bss.TopicPostUpdate {
	status := constants.BeatmapStatusWIP
	if request.Complete {
		status = constants.BeatmapStatusPending
	}

	return bss.TopicPostUpdate{
		Subject:   request.Subject,
		Message:   request.Message,
		Status:    status,
		Subscribe: request.Subscribe,
		Bump:      true,
	}
}

// /web/osu-get-beatmap-topic.php -> Returns the forum topic contents assigned to a beatmapset.
func BeatmapSubmissionGetTopic(ctx *server.Context) {
	request, err := NewOsz2GetTopicRequest(ctx)
	if err != nil {
		renderBssError(ctx, fmt.Errorf("%w: %v", bss.ErrInvalidSubmissionRequest, err))
		return
	}

	user, err := ctx.AuthenticateUser(
		request.Username,
		request.Password,
		false, // doesn't require a bancho presence
	)
	if errors.Is(err, server.ErrUserNotFound) || errors.Is(err, server.ErrInvalidPassword) {
		renderBssError(ctx, bss.ErrAuthenticationFailed)
		return
	}
	if err != nil {
		renderBssError(ctx, fmt.Errorf("authenticate user: %w", err))
		return
	}

	beatmapset, err := ctx.State.Beatmapsets.ById(request.SetId)
	if err != nil {
		renderBssError(ctx, fmt.Errorf("fetch beatmapset %d: %w", request.SetId, err))
		return
	}
	if beatmapset == nil {
		ctx.Logger.Warn("Failed to fetch beatmap topic: beatmapset not found", "set_id", request.SetId)
		renderBssError(ctx, bss.ErrBeatmapsetNotFound)
		return
	}
	if beatmapset.CreatorId == nil {
		ctx.Logger.Warn("Failed to fetch beatmap topic: beatmapset creator not set", "set_id", request.SetId)
		renderBssError(ctx, bss.ErrBeatmapsetWrongServer)
		return
	}
	if *beatmapset.CreatorId != user.Id {
		ctx.Logger.Warn("Failed to fetch beatmap topic: beatmapset creator mismatch", "set_id", request.SetId)
		renderBssError(ctx, bss.ErrBeatmapsetAccessNotAllowed)
		return
	}

	contents, err := bss.GetTopic(ctx.State.Repositories, beatmapset)
	if err != nil {
		renderBssError(ctx, err)
		return
	}

	response := strings.Join([]string{
		"0",
		strconv.Itoa(contents.Id),
		contents.Subject,
		contents.Message,
	}, "\x03")
	ctx.RenderText(http.StatusOK, response)
}

// /web/osu-osz2-bmsubmit-post.php -> Creates / updates the forum topic & sets beatmapset status.
func BeatmapSubmissionOsz2Post(ctx *server.Context) {
	request, err := NewOsz2PostTopicRequest(ctx)
	if err != nil {
		renderStatusError(ctx, fmt.Errorf("%w: %v", bss.ErrInvalidSubmissionRequest, err))
		return
	}

	user, err := ctx.AuthenticateUser(
		request.Username,
		request.Password,
		false, // does not require bancho presence
	)
	if errors.Is(err, server.ErrUserNotFound) || errors.Is(err, server.ErrInvalidPassword) {
		renderStatusError(ctx, bss.ErrAuthenticationFailed)
		return
	}
	if err != nil {
		renderStatusError(ctx, fmt.Errorf("authenticate user: %w", err))
		return
	}

	userPermissions, err := ctx.State.Permissions.Resolve(user.Id)
	if err != nil {
		renderStatusError(ctx, fmt.Errorf("resolve permissions for user %d: %w", user.Id, err))
		return
	}

	submission := bss.NewSubmissionContext(ctx.Context())
	submission.User = user

	if err := submission.CheckSubmissionEligibility(
		ctx.State.Config.BeatmapSubmissionEnabled,
		userPermissions,
	); err != nil {
		renderStatusError(ctx, err)
		return
	}

	if !userPermissions.Has("forum.posts.create") {
		renderStatusError(ctx, bss.ErrUploadPermissionDenied)
	}
	if !userPermissions.Has("forum.posts.edit") {
		renderStatusError(ctx, bss.ErrUploadPermissionDenied)
	}
	var topicId int

	err = ctx.State.DatabaseTransaction(func(transaction *state.Repositories) error {
		if err := submission.PrepareBeatmapset(transaction, request.SetId); err != nil {
			return err
		}
		if err := submission.RequireBeatmapsetOwner(); err != nil {
			return err
		}
		if err := submission.CheckBeatmapsetEligibility(); err != nil {
			return err
		}

		topic, err := bss.PostTopic(
			transaction,
			submission.Beatmapset,
			request.TopicUpdate(),
		)
		if err != nil {
			return err
		}
		topicId = topic.Id
		return nil
	})
	if err != nil {
		renderStatusError(ctx, err)
		return
	}

	ctx.Logger.Info(
		"Beatmap forum post updated",
		"user_id", user.Id,
		"set_id", request.SetId,
		"topic_id", topicId,
		"status", request.TopicUpdate().Status,
	)
	ctx.RenderText(http.StatusOK, strconv.Itoa(topicId))
}

func renderStatusError(ctx *server.Context, err error) {
	status := bss.StatusCodeFor(err)
	expected := status != http.StatusInternalServerError

	if expected {
		ctx.Logger.Warn("Rejected osz2 topic submission", "error", err)
	} else {
		ctx.Logger.Error("Failed to process osz2 topic submission", "error", err)
	}
	ctx.Response.WriteHeader(status)
}
