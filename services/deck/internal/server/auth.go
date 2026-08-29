package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/osuTitanic/titanic/internal/authentication"
	"github.com/osuTitanic/titanic/internal/schemas"
)

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrInvalidPassword        = errors.New("invalid password")
	ErrBanchoPresenceNotFound = errors.New("bancho presence not found")
)

// AuthenticateUser authenticates a user by their username and password.
// If `requireBanchoPresence` is true, the user must also be present on Bancho.
// Use ErrUserNotFound, ErrInvalidPassword, and ErrBanchoPresenceNotFound to check for specific errors.
func (ctx *Context) AuthenticateUser(
	username string,
	password string,
	requireBanchoPresence bool,
) (*schemas.User, error) {
	user, err := ctx.State.Users.ByName(username)
	if err != nil {
		return nil, fmt.Errorf("authenticate user %q: %w", username, err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	if !authentication.VerifyPasswordHashFromMd5(password, user.Bcrypt) {
		return nil, ErrInvalidPassword
	}

	if !requireBanchoPresence {
		return user, nil
	}

	online, err := ctx.State.Redis.Exists(
		ctx.Request.Context(),
		"bancho:status:"+strconv.Itoa(user.Id),
	).Result()
	if err != nil {
		return nil, fmt.Errorf("check bancho presence for user %d: %w", user.Id, err)
	}
	if online == 0 {
		return nil, ErrBanchoPresenceNotFound
	}

	return user, nil
}

// AuthenticateUserFromQuery authenticates a user using query parameters for username and password.
func (ctx *Context) AuthenticateUserFromQuery(
	usernameKey string,
	passwordKey string,
	requireBanchoPresence bool,
) (*schemas.User, error) {
	return ctx.AuthenticateUser(
		ctx.QueryValue(usernameKey),
		ctx.QueryValue(passwordKey),
		requireBanchoPresence,
	)
}

// HandleUserAuthenticationSimple authenticates a user and handles the response in case of failure.
func (ctx *Context) HandleUserAuthenticationSimple(
	usernameKey string,
	passwordKey string,
	requireBanchoPresence bool,
) (*schemas.User, bool) {
	user, err := ctx.AuthenticateUserFromQuery(
		usernameKey,
		passwordKey,
		requireBanchoPresence,
	)

	switch {
	case errors.Is(err, ErrUserNotFound),
		errors.Is(err, ErrInvalidPassword),
		errors.Is(err, ErrBanchoPresenceNotFound):
		ctx.Response.WriteHeader(http.StatusUnauthorized)
	case err != nil:
		ctx.Logger.Error("Failed to authenticate user", "error", err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
	default:
		return user, true
	}
	return nil, false
}
