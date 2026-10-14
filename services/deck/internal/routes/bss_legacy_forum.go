package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/services/deck/internal/bss"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

type PreOsz2PostTopicRequest struct {
	Username  string
	Password  string
	SetId     int
	Subject   string
	Message   string
	Complete  bool
	Subscribe bool
	Bump      bool
}

func NewPreOsz2PostTopicRequest(ctx *server.Context) (request PreOsz2PostTopicRequest, err error) {
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
	bump, _ := ctx.FormValueInt("bumprequest")
	request.Complete = complete == 1
	request.Subscribe = subscribe == 1
	request.Bump = bump == 1
	return request, nil
}

func (request PreOsz2PostTopicRequest) TopicUpdate() bss.TopicPostUpdate {
	status := constants.BeatmapStatusWIP
	if request.Complete {
		status = constants.BeatmapStatusPending
	}

	return bss.TopicPostUpdate{
		Subject:   request.Subject,
		Message:   request.Message,
		Status:    status,
		Subscribe: request.Subscribe,
		Bump:      request.Bump,
	}
}

// /web/osu-bmsubmit-post[2-3].php -> Creates / updates the forum topic & sets beatmapset status.
func BeatmapSubmissionPreOsz2Post(ctx *server.Context) {
	request, err := NewPreOsz2PostTopicRequest(ctx)
	if err != nil {
		renderTopicPostError(ctx, fmt.Errorf("%w: %v", bss.ErrInvalidSubmissionRequest, err))
		return
	}

	submission, topicId, err := postBeatmapTopic(
		ctx,
		request.Username,
		request.Password,
		request.SetId,
		request.TopicUpdate(),
	)
	if err != nil {
		renderTopicPostError(ctx, err)
		return
	}

	ctx.Logger.Info(
		"Beatmap forum post updated",
		"user_id", submission.User.Id,
		"set_id", request.SetId,
		"topic_id", topicId,
		"status", request.TopicUpdate().Status,
		"bump", request.Bump,
	)
	ctx.RenderText(http.StatusOK, strconv.Itoa(topicId))
}
