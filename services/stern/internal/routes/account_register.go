package routes

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/osuTitanic/titanic/internal/activity"
	"github.com/osuTitanic/titanic/internal/authentication"
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/internal/discord"
	"github.com/osuTitanic/titanic/internal/schemas"
	"github.com/osuTitanic/titanic/internal/state"
	"github.com/osuTitanic/titanic/services/stern/internal/server"
	"github.com/osuTitanic/titanic/services/stern/internal/templates"
	"github.com/redis/go-redis/v9"
)

type registrationRequest struct {
	Username       string
	SafeName       string
	Email          string
	HashedPassword string
	Country        string
}

type registrationResult struct {
	User         *schemas.User
	Verification *schemas.Verification
}

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

func AccountRegister(ctx *server.Context) {
	if ctx.IsAuthenticated() {
		ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/u/%d", ctx.CurrentUser.Id))
		return
	}

	if err := ctx.Request.ParseForm(); err != nil {
		ctx.Logger.Warn("Failed to parse registration form", "error", err)
		RenderRegisterPage(ctx, "Failed to process your request. Please try again!")
		return
	}

	// Validate username, email & password inputs
	username := strings.TrimSpace(ctx.Request.FormValue("username"))
	email := strings.ToLower(strings.TrimSpace(ctx.Request.FormValue("email")))
	password := ctx.Request.FormValue("password")

	validationError, err := validateRegistrationEmail(ctx, email)
	if err != nil {
		ctx.Logger.Error("Failed to validate registration email", "error", err)
		RenderRegisterPage(ctx, "Failed to process your request. Please try again!")
		return
	}
	if validationError != "" {
		RenderRegisterPage(ctx, "Please enter a valid email!")
		return
	}

	validationError, err = validateRegistrationUsername(ctx, username)
	if err != nil {
		ctx.Logger.Error("Failed to validate registration username", "error", err)
		RenderRegisterPage(ctx, "Failed to process your request. Please try again!")
		return
	}
	if validationError != "" {
		RenderRegisterPage(ctx, "Please enter a valid username!")
		return
	}

	if password == "" {
		RenderRegisterPage(ctx, "Please enter a valid password!")
		return
	}
	if len(password) <= 7 {
		RenderRegisterPage(ctx, "Please enter a password with at least 8 characters!")
		return
	}

	// Check if there have been too many registrations from this ip in the last 24 hours
	tooManyRegistrations, err := hasTooManyRegistrations(ctx)
	if err != nil {
		ctx.Logger.Error("Failed to resolve registration count", "error", err)
		InternalServerError(ctx)
		return
	}
	if tooManyRegistrations {
		RenderRegisterPage(ctx, "There have been too many registrations from this ip. Please try again later!")
		return
	}

	// Validate recaptcha response, if enabled
	if ctx.State.Config.RecaptchaSecretKey != "" && ctx.State.Config.RecaptchaSiteKey != "" {
		clientResponse := strings.TrimSpace(ctx.Request.FormValue("recaptcha_response"))
		if clientResponse == "" {
			RenderRegisterPage(ctx, "Invalid captcha response!")
			return
		}

		ok, err := server.VerifyRecaptchaResponse(ctx, clientResponse)
		if err != nil {
			ctx.Logger.Warn("Failed to verify registration captcha", "error", err)
			RenderRegisterPage(ctx, "Failed to verify captcha response!")
			return
		}
		if !ok {
			RenderRegisterPage(ctx, "Captcha verification failed!")
			return
		}
	}

	hashedPassword, err := authentication.CreatePasswordHash(password)
	if err != nil {
		ctx.Logger.Error("Failed to hash registration password", "error", err)
		RenderRegisterPage(ctx, "An error occured on the server side. Please try again!")
		return
	}

	input := registrationRequest{
		Username:       username,
		SafeName:       schemas.ResolveSafeName(username),
		Country:        ctx.Country(),
		Email:          email,
		HashedPassword: hashedPassword,
	}

	result, err := performRegistration(ctx, input)
	if err != nil {
		ctx.Logger.Error("Failed to create registration", "username", username, "email", email, "error", err)
		RenderRegisterPage(ctx, "An error occured on the server side. Please try again!")
		return
	}
	notifyOfficerAboutRegistration(ctx, result.User)
	broadcastRegistrationActivity(ctx, result.User)

	if err := sendWelcomeEmail(ctx, result.Verification); err != nil {
		ctx.Logger.Error("Failed to send registration verification email", "user_id", result.User.Id, "verification_id", result.Verification.Id, "error", err)
		RenderRegisterPage(ctx, "Failed to send verification email. Please try again later!")
		return
	}

	if err := recordRegistrationForIP(ctx); err != nil {
		ctx.Logger.Warn("Failed to record registration count", "user_id", result.User.Id, "error", err)
	}

	if result.Verification != nil || result.User.Activated {
		// Redirect to account verification info page
		ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/account/verification?id=%d", result.Verification.Id))
		return
	}

	// No verification / User was already activated automatically
	ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/u/%d", result.User.Id))
}

func AccountRegisterCheck(ctx *server.Context) {
	fieldType := ctx.QueryValue("type")
	value := ctx.QueryValue("value")
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

func performRegistration(ctx *server.Context, input registrationRequest) (result *registrationResult, err error) {
	result = &registrationResult{}
	err = ctx.State.DatabaseTransaction(func(repositories *state.Repositories) error {
		user := &schemas.User{
			Name:      input.Username,
			SafeName:  input.SafeName,
			Email:     input.Email,
			Bcrypt:    input.HashedPassword,
			Country:   input.Country,
			Activated: false, // TODO: Auto-activate if no email provider is given
		}
		if err := repositories.Users.Create(user); err != nil {
			return err
		}

		playerGroup := &schemas.GroupEntry{UserId: user.Id, GroupId: constants.GroupPlayers}
		supporterGroup := &schemas.GroupEntry{UserId: user.Id, GroupId: constants.GroupSupporter}

		if err := repositories.Groups.CreateEntry(playerGroup); err != nil {
			return err
		}
		if err := repositories.Groups.CreateEntry(supporterGroup); err != nil {
			return err
		}

		err := repositories.Notifications.Create(&schemas.Notification{
			UserId:  user.Id,
			Type:    constants.NotificationTypeWelcome,
			Header:  constants.WelcomeNotificationHeader,
			Content: fmt.Sprintf(constants.WelcomeNotificationContent, ctx.State.Config.OsuBaseUrl()),
			Link:    "/download",
		})
		if err != nil {
			return err
		}

		token, err := generateVerificationToken()
		if err != nil {
			return err
		}

		verification, err := repositories.Verifications.CreateForUser(
			user.Id,
			constants.VerificationTypeActivation,
			token,
			time.Now(),
		)
		if err != nil {
			return err
		}

		result.User = user
		result.Verification = verification
		result.Verification.User = user
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func notifyOfficerAboutRegistration(ctx *server.Context, user *schemas.User) {
	if ctx.State == nil {
		return
	}

	title := "New registration"
	description := fmt.Sprintf(
		"[%s](%s/u/%d) registered an account.",
		user.Name,
		strings.TrimRight(ctx.State.Config.OsuBaseUrl(), "/"),
		user.Id,
	)
	timestamp := time.Now()
	color := 0x66CCFF

	embed := discord.Embed{
		Title:       &title,
		Description: &description,
		Color:       &color,
		Timestamp:   &timestamp,
	}
	embed.AddField("User ID", fmt.Sprintf("`%d`", user.Id), true)
	embed.AddField("Username", fmt.Sprintf("`%s`", user.Name), true)
	embed.AddField("Country", fmt.Sprintf("`%s`", user.Country), true)
	embed.AddField("IP", fmt.Sprintf("||`%s`||", ctx.IP()), true)

	if err := ctx.State.Officer.Call(discord.OfficerTagRegistration, "", embed); err != nil {
		ctx.Logger.Warn("Failed to send registration officer notification", "user_id", user.Id, "error", err)
	}
}

func broadcastRegistrationActivity(ctx *server.Context, user *schemas.User) {
	err := activity.Submit(
		ctx.State, user.Id, nil,
		constants.ActivityUserRegistration,
		map[string]any{"username": user.Name},
		false, // should not be broadcasted, i.e. only shown in activity websocket
		true,  // should not be stored in db
	)
	if err != nil {
		ctx.Logger.Warn("Failed to record registration activity", "user_id", user.Id, "error", err)
	}
}

func hasTooManyRegistrations(ctx *server.Context) (bool, error) {
	key := "registrations:" + ctx.IP()
	registrations, err := ctx.State.Redis.Get(ctx.Request.Context(), key).Int()
	if err == nil {
		return registrations > 2, nil
	}
	if err == redis.Nil {
		return false, nil
	}
	return false, err
}

func recordRegistrationForIP(ctx *server.Context) error {
	key := "registrations:" + ctx.IP()
	if err := ctx.State.Redis.Incr(ctx.Request.Context(), key).Err(); err != nil {
		return err
	}
	return ctx.State.Redis.Expire(ctx.Request.Context(), key, 24*time.Hour).Err()
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

	user, err := ctx.State.Users.ByEmail(strings.ToLower(email))
	if err != nil {
		return "", err
	}
	if user != nil {
		return "This email address is already in use.", nil
	}
	return "", nil
}

func registrationUserExists(ctx *server.Context, username string) (bool, error) {
	user, err := ctx.State.Users.ByNameCaseInsensitive(username)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}

func registrationSafeNameExists(ctx *server.Context, safeName string) (bool, error) {
	user, err := ctx.State.Users.BySafeName(safeName)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}

func registrationReservedNameExists(ctx *server.Context, username string) (bool, error) {
	name, err := ctx.State.Names.ByReservedNameCaseInsensitive(username)
	if err != nil {
		return false, err
	}
	return name != nil, nil
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
