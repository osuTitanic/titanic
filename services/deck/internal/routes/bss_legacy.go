package routes

import (
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/services/deck/internal/bss"
)

// TODO: /web/osu-bmsubmit-getid[5-2].php
// TODO: /web/osu-bmsubmit-upload.php
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

func (request *PreOsz2GetIdRequest) Size() int64 {
	if request.FileHeader == nil {
		return -1
	}
	return request.FileHeader.Size
}

func (request *PreOsz2GetIdRequest) Close() {
	if request.File == nil {
		return
	}
	request.File.Close()
	request.File = nil
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
