package bss

import (
	"errors"
	"fmt"
	"net/http"
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
	ErrReservedSetNotFound   = errors.New("bss: reserved beatmapset not found")
	ErrBeatmapsetWrongServer = errors.New("bss: beatmapset is hosted on bancho")
	ErrBeatmapsetRanked      = errors.New("bss: beatmapset is approved")
	ErrBeatmapsetGraveyarded = errors.New("bss: beatmapset is graveyarded")

	ErrBeatmapsetAccessNotAllowed = errors.New("bss: beatmapset access not allowed")
	ErrBeatmapAccessNotAllowed    = errors.New("bss: beatmap access not allowed")
	ErrBeatmapCreationNotAllowed  = errors.New("bss: beatmap creation not allowed")
	ErrBeatmapRemovalNotAllowed   = errors.New("bss: beatmap removal not allowed")
	ErrResourceAccessNotAllowed   = errors.New("bss: resource access not allowed")
	ErrNoUploadSlots              = errors.New("bss: no remaining upload slots")
	ErrBeatmapOwnershipNotAllowed = errors.New("bss: beatmap creator is not allowed")
	ErrDuplicateBeatmap           = errors.New("bss: duplicate beatmap")
	ErrBeatmapTooShort            = errors.New("bss: beatmap is too short")

	ErrPackageInvalid  = errors.New("bss: invalid package")
	ErrPackageTooLarge = errors.New("bss: package is too large")
	ErrPackageNotFound = errors.New("bss: package not found")
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
	case errors.Is(err, ErrBeatmapsetAccessNotAllowed),
		errors.Is(err, ErrBeatmapAccessNotAllowed),
		errors.Is(err, ErrResourceAccessNotAllowed),
		errors.Is(err, ErrBeatmapsetWrongServer),
		errors.Is(err, ErrBeatmapOwnershipNotAllowed):
		return NewErrorResponse(ErrorCodeOwnership), true
	case errors.Is(err, ErrInvalidSubmissionRequest), errors.Is(err, ErrInvalidSubmissionTicket):
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
	case errors.Is(err, ErrReservedSetNotFound):
		return NewErrorResponseMessage("The reserved beatmapset no longer exists. Please try again!"), true
	case errors.Is(err, ErrBeatmapRemovalNotAllowed):
		return NewErrorResponseMessage("Please ask the owner of this beatmapset to delete any difficulties."), true
	case errors.Is(err, ErrDuplicateBeatmap):
		return NewErrorResponseMessage("It seems like one of your beatmaps was already uploaded by someone else. Please try again!"), true
	case errors.Is(err, ErrPackageNotFound):
		return NewErrorResponseMessage("Something went wrong while processing your submission. Please try again!"), true
	case errors.Is(err, ErrBeatmapTooShort):
		return NewErrorResponseMessage("Your beatmap is too short. Please try to make it longer and try again!"), true
	case errors.Is(err, ErrPackageTooLarge):
		return NewErrorResponseMessage("Your beatmap is too big. Try to reduce its filesize and try again!"), true
	case errors.Is(err, ErrPackageInvalid):
		return NewErrorResponseMessage("Something went wrong while processing your beatmap. Please try again!"), true
	default:
		return ErrorResponse{}, false
	}
}

func StatusCodeFor(err error) int {
	switch {
	case errors.Is(err, ErrAuthenticationFailed),
		errors.Is(err, ErrUserNotAuthenticated),
		errors.Is(err, ErrUserNotActivated),
		errors.Is(err, ErrUserRestricted),
		errors.Is(err, ErrUserSilenced),
		errors.Is(err, ErrUserIsBot),
		errors.Is(err, ErrUploadPermissionDenied),
		errors.Is(err, ErrBeatmapsetAccessNotAllowed),
		errors.Is(err, ErrBeatmapsetWrongServer),
		errors.Is(err, ErrBeatmapsetRanked),
		errors.Is(err, ErrBeatmapsetGraveyarded):
		return http.StatusForbidden
	case errors.Is(err, ErrInvalidSubmissionRequest):
		return http.StatusBadRequest
	case errors.Is(err, ErrBeatmapsetNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrSubmissionDisabled):
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

func FormatOsz2Error(err error) (string, bool) {
	response, expected := ErrorResponseFor(err)
	if !expected {
		response = GenericErrorResponse()
	}
	return fmt.Sprintf("%d\n%s", response.Code, response.Message), expected
}

func FormatPreOsz2Error(err error) (string, bool) {
	switch {
	case errors.Is(err, ErrBeatmapsetRanked):
		return "ranked", true
	case errors.Is(err, ErrBeatmapsetGraveyarded):
		return "graveyarded", true
	}

	response, expected := ErrorResponseFor(err)
	if response.Message != "" {
		return response.Message, expected
	}

	switch response.Code {
	case ErrorCodeOwnership:
		return "The beatmap you're trying to submit isn't owned by you.", expected
	case ErrorCodeUnavailable:
		return "The beatmap you're trying to submit is no longer available.", expected
	case ErrorCodeQuotaExceeded:
		return "You have reached your maximum amount of beatmaps you can upload.", expected
	default:
		return "An error occurred while processing your beatmap. Please try again!", false
	}
}
