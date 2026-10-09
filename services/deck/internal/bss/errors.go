package bss

import (
	"errors"
	"fmt"
)

type ErrorCode int

const (
	ErrorCodeOwnership ErrorCode = iota + 1
	ErrorCodeUnavailable
	ErrorCodeRanked
	ErrorCodeGraveyarded
	ErrorCodeCustom
	ErrorCodeQuotaExceeded
)

type ErrorResponse struct {
	Code    ErrorCode
	Message string
}

var (
	ErrInvalidSubmissionRequest = errors.New("bss: invalid submission request")
	ErrAuthenticationFailed     = errors.New("bss: authentication failed")
	ErrSubmissionDisabled       = errors.New("bss: beatmap submission disabled")
	ErrUserNotAuthenticated     = errors.New("bss: user not authenticated")
	ErrUserNotActivated         = errors.New("bss: user not activated")
	ErrUserRestricted           = errors.New("bss: user restricted")
	ErrUserSilenced             = errors.New("bss: user silenced")
	ErrUserIsBot                = errors.New("bss: bot users cannot submit beatmaps")
	ErrUploadPermissionDenied   = errors.New("bss: beatmap upload permission denied")

	ErrBeatmapsetNotPrepared = errors.New("bss: beatmapset not prepared")
	ErrBeatmapsetNotFound    = errors.New("bss: beatmapset not found")
	ErrBeatmapsetWrongServer = errors.New("bss: beatmapset is hosted on bancho")
	ErrBeatmapsetRanked      = errors.New("bss: beatmapset is approved")
	ErrBeatmapsetGraveyarded = errors.New("bss: beatmapset is graveyarded")

	ErrBeatmapsetAccessNotAllowed = errors.New("bss: beatmapset access not allowed")
	ErrBeatmapAccessNotAllowed    = errors.New("bss: beatmap access not allowed")
	ErrBeatmapCreationNotAllowed  = errors.New("bss: beatmap creation not allowed")
	ErrBeatmapRemovalNotAllowed   = errors.New("bss: beatmap removal not allowed")
	ErrResourceAccessNotAllowed   = errors.New("bss: resource access not allowed")
	ErrNoUploadSlots              = errors.New("bss: no remaining upload slots")
)

func NewErrorResponse(code ErrorCode) ErrorResponse {
	return ErrorResponse{Code: code}
}

func NewErrorResponseMessage(message string) ErrorResponse {
	return ErrorResponse{Code: ErrorCodeCustom, Message: message}
}

func GenericErrorResponse() ErrorResponse {
	return NewErrorResponseMessage("A server error occurred. Please try again!")
}

func ErrorResponseFor(err error) (ErrorResponse, bool) {
	switch {
	case errors.Is(err, ErrBeatmapsetRanked):
		return NewErrorResponse(ErrorCodeRanked), true
	case errors.Is(err, ErrBeatmapsetGraveyarded):
		return NewErrorResponse(ErrorCodeGraveyarded), true
	case errors.Is(err, ErrNoUploadSlots):
		return NewErrorResponse(ErrorCodeQuotaExceeded), true
	case errors.Is(err, ErrBeatmapsetAccessNotAllowed), errors.Is(err, ErrBeatmapAccessNotAllowed), errors.Is(err, ErrResourceAccessNotAllowed), errors.Is(err, ErrBeatmapsetWrongServer):
		return NewErrorResponse(ErrorCodeOwnership), true
	case errors.Is(err, ErrInvalidSubmissionRequest):
		return NewErrorResponseMessage("Invalid beatmap submission request."), true
	case errors.Is(err, ErrAuthenticationFailed), errors.Is(err, ErrUserNotAuthenticated):
		return NewErrorResponseMessage("Authentication failed. Please check your username and password and try again!"), true
	case errors.Is(err, ErrSubmissionDisabled):
		return NewErrorResponseMessage("The beatmap submission system is currently disabled. Please try again later!"), true
	case errors.Is(err, ErrUserNotActivated):
		return NewErrorResponseMessage("Please activate your account before uploading beatmaps."), true
	case errors.Is(err, ErrUserRestricted):
		return NewErrorResponseMessage("You are banned. Please contact support if you believe this is a mistake."), true
	case errors.Is(err, ErrUserSilenced):
		return NewErrorResponseMessage("You are not allowed to upload beatmaps while silenced."), true
	case errors.Is(err, ErrUserIsBot):
		return NewErrorResponseMessage("Bot accounts cannot upload beatmaps."), true
	case errors.Is(err, ErrUploadPermissionDenied):
		return NewErrorResponseMessage("You do not have permission to upload beatmaps."), true
	case errors.Is(err, ErrBeatmapCreationNotAllowed):
		return NewErrorResponseMessage("You are not allowed to add difficulties to this beatmapset."), true
	case errors.Is(err, ErrBeatmapRemovalNotAllowed):
		return NewErrorResponseMessage("Please ask the owner of this beatmapset to delete any difficulties."), true
	default:
		return ErrorResponse{}, false
	}
}

func FormatOsz2ErrorResponse(err error) (string, bool) {
	response, expected := ErrorResponseFor(err)
	if !expected {
		response = GenericErrorResponse()
	}
	return fmt.Sprintf("%d\n%s", response.Code, response.Message), expected
}
