package server

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/osuTitanic/titanic/internal/authentication"
	"github.com/osuTitanic/titanic/internal/schemas"
)

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrInvalidPassword        = errors.New("invalid password")
	ErrBanchoPresenceNotFound = errors.New("bancho presence not found")
)

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
