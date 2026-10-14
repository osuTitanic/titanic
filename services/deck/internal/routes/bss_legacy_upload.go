package routes

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"path"
	"strings"

	"github.com/osuTitanic/titanic/services/deck/internal/bss"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

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
