package routes

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/osuTitanic/titanic-go/internal/constants"
	"github.com/osuTitanic/titanic-go/internal/email"
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/internal/state"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func PasswordResetPage(ctx *server.Context) {
	if ctx.IsAuthenticated() {
		ctx.Redirect(http.StatusSeeOther, "/")
		return
	}
	RenderResetPage(ctx, "")
}

func RenderResetPage(ctx *server.Context, errorMessage string) {
	view := templates.ResetView{
		DefaultView:  buildDefaultView(ctx),
		ErrorMessage: errorMessage,
	}
	ctx.RenderTemplate(http.StatusOK, "pages/account/reset", view)
}

func PasswordReset(ctx *server.Context) {
	if ctx.IsAuthenticated() {
		ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/u/%d", ctx.CurrentUser.Id))
		return
	}

	if err := ctx.Request.ParseForm(); err != nil {
		ctx.Logger.Warn("Failed to parse password reset form", "error", err)
		RenderResetPage(ctx, "Failed to process your request. Please try again!")
		return
	}

	if !ctx.State.Config.EmailsEnabled() {
		ctx.Logger.Warn("Password reset requested but emails are disabled")
		RenderResetPage(ctx, "Password resets are not enabled at the moment. Please contact an administrator!")
		return
	}

	emailAddress := strings.TrimSpace(ctx.Request.FormValue("email"))
	emailAddressLower := strings.ToLower(emailAddress)
	if emailAddress == "" {
		RenderResetPage(ctx, "Please enter a valid email!")
		return
	}

	user, err := ctx.State.Users.ByEmail(emailAddressLower)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			RenderResetPage(ctx, "We could not find any user with that email address.")
			return
		}
		ctx.Logger.Error("Failed to fetch user for password reset", "email", emailAddress, "error", err)
		InternalServerError(ctx)
		return
	}

	// Prevent users from spamming password reset requests
	locked, err := hasPasswordResetLock(ctx, user.Id)
	if err != nil {
		ctx.Logger.Error("Failed to resolve password reset lock", "user_id", user.Id, "error", err)
		InternalServerError(ctx)
		return
	}
	if locked {
		RenderResetPage(ctx, "You have already requested a password reset recently. Please check your emails, or try again in a few hours!")
		return
	}

	ctx.Logger.Info("Sending verification email for resetting password...", "user_id", user.Id)

	verification, err := createPasswordResetVerification(ctx, user)
	if err != nil {
		ctx.Logger.Error("Failed to create password reset verification", "user_id", user.Id, "error", err)
		InternalServerError(ctx)
		return
	}

	if err := sendPasswordResetEmail(ctx, verification); err != nil {
		ctx.Logger.Error("Failed to send password reset email", "user_id", user.Id, "verification_id", verification.Id, "error", err)
		RenderResetPage(ctx, "Failed to send password reset email. Please try again later!")
		return
	}

	// Set a lock for the user to prevent spamming
	if err := setPasswordResetLock(ctx, user.Id); err != nil {
		ctx.Logger.Warn("Failed to set password reset lock", "user_id", user.Id, "error", err)
	}

	ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/account/verification?id=%d", verification.Id))
}

func createPasswordResetVerification(ctx *server.Context, user *schemas.User) (*schemas.Verification, error) {
	var verification *schemas.Verification

	// Use transaction to ensure that the verification is created atomically
	err := ctx.State.DatabaseTransaction(func(repos *state.Repositories) error {
		token, err := generateVerificationToken()
		if err != nil {
			return err
		}

		verification, err = repos.Verifications.CreateForUser(
			user.Id,
			constants.VerificationTypePassword,
			token,
			time.Now(),
		)
		return err
	})
	if err != nil {
		return nil, err
	}

	verification.User = user
	return verification, nil
}

func sendPasswordResetEmail(ctx *server.Context, verification *schemas.Verification) error {
	if verification.User == nil {
		return errors.New("password reset email: missing verification user")
	}
	osuBaseUrl := ctx.State.Config.OsuBaseUrl()

	body := fmt.Sprintf(
		constants.EmailPasswordResetBody,
		verification.User.Name,
		osuBaseUrl,
		verification.Id,
		verification.Token,
		osuBaseUrl,
	)
	return ctx.State.Email.Send(&email.Message{
		To:       []string{verification.User.Email},
		Subject:  constants.EmailPasswordResetSubject,
		TextBody: body,
	})
}

func hasPasswordResetLock(ctx *server.Context, userId int) (bool, error) {
	value, err := ctx.State.Redis.Get(ctx.Request.Context(), fmt.Sprintf("reset_lock:%d", userId)).Int()
	if err == nil {
		return value != 0, nil
	}
	if err == redis.Nil {
		return false, nil
	}
	return false, err
}

func setPasswordResetLock(ctx *server.Context, userId int) error {
	return ctx.State.Redis.Set(ctx.Request.Context(), fmt.Sprintf("reset_lock:%d", userId), 1, 12*time.Hour).Err()
}
