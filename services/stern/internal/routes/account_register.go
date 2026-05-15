package routes

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/osuTitanic/titanic-go/internal/constants"
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
	"gorm.io/gorm"
)

func AccountRegisterPage(ctx *server.Context) {
	if ctx.IsAuthenticated() {
		ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/u/%d", ctx.CurrentUser.Id))
		return
	}
	RenderRegisterPage(ctx, "")
}

func RenderRegisterPage(ctx *server.Context, errorMessage string) {
	view := templates.RegisterView{
		DefaultView:      buildDefaultView(ctx),
		ErrorMessage:     errorMessage,
		RecaptchaEnabled: ctx.State.Config.RecaptchaSiteKey != "",
		RecaptchaSiteKey: ctx.State.Config.RecaptchaSiteKey,
	}
	ctx.RenderTemplate(http.StatusOK, "pages/account/register", view)
}

func AccountRegisterCheck(ctx *server.Context) {
	fieldType := ctx.Request.URL.Query().Get("type")
	value := ctx.Request.URL.Query().Get("value")
	if fieldType == "" || value == "" {
		writePlainText(ctx, http.StatusOK, "")
		return
	}

	var (
		validationError string
		err             error
	)
	switch fieldType {
	case "username":
		validationError, err = validateRegistrationUsername(ctx, value)
	case "email":
		validationError, err = validateRegistrationEmail(ctx, value)
	default:
		writePlainText(ctx, http.StatusOK, "")
		return
	}
	if err != nil {
		ctx.Logger.Error("Failed to validate registration field", "type", fieldType, "error", err)
		writePlainText(ctx, http.StatusInternalServerError, "Could not verify this field. Please try something else!")
		return
	}

	writePlainText(ctx, http.StatusOK, validationError)
}

func validateRegistrationUsername(ctx *server.Context, username string) (string, error) {
	username = strings.TrimSpace(username)
	if len(username) < 3 {
		return "Your username is too short.", nil
	}
	if len(username) > 15 {
		return "Your username is too long.", nil
	}
	if !constants.Username.MatchString(username) {
		return "Your username contains invalid characters.", nil
	}

	lowerUsername := strings.ToLower(username)
	for _, word := range constants.DisallowedUsernameSubstrings {
		if strings.Contains(lowerUsername, word) {
			return "Your username contains offensive words.", nil
		}
	}
	if strings.HasPrefix(lowerUsername, "deleteduser") {
		return "This username is not allowed.", nil
	}
	if strings.HasSuffix(lowerUsername, "_old") {
		return "This username is not allowed.", nil
	}

	if exists, err := registrationUserExists(ctx, username); err != nil {
		return "", err
	} else if exists {
		return "This username is already in use!", nil
	}

	safeName := schemas.ResolveSafeName(username)
	if exists, err := registrationSafeNameExists(ctx, safeName); err != nil {
		return "", err
	} else if exists {
		return "This username is already in use!", nil
	}

	if exists, err := registrationReservedNameExists(ctx, username); err != nil {
		return "", err
	} else if exists {
		return "This username is already in use!", nil
	}

	return "", nil
}

func validateRegistrationEmail(ctx *server.Context, email string) (string, error) {
	if !constants.Email.MatchString(email) {
		return "Please enter a valid email address!", nil
	}

	_, err := ctx.State.Users.ByEmail(strings.ToLower(email))
	if err == nil {
		return "This email address is already in use.", nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	}
	return "", err
}

func registrationUserExists(ctx *server.Context, username string) (bool, error) {
	_, err := ctx.State.Users.ByNameCaseInsensitive(username)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}

func registrationSafeNameExists(ctx *server.Context, safeName string) (bool, error) {
	_, err := ctx.State.Users.BySafeName(safeName)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}

func registrationReservedNameExists(ctx *server.Context, username string) (bool, error) {
	_, err := ctx.State.Names.ByReservedNameCaseInsensitive(username)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}

func writePlainText(ctx *server.Context, status int, body string) {
	ctx.Response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	ctx.Response.WriteHeader(status)
	if body == "" {
		return
	}
	if _, err := ctx.Response.Write([]byte(body)); err != nil {
		ctx.Logger.Error("Failed to write plain text response", "error", err)
	}
}
