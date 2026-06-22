package routes

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/osuTitanic/titanic-go/internal/authentication"
	"github.com/osuTitanic/titanic-go/internal/constants"
	"github.com/osuTitanic/titanic-go/internal/email"
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/internal/state"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
	"github.com/redis/go-redis/v9"
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

	// If a verification token is present, the user has submitted a new
	// password from the verification page
	if token := ctx.Request.FormValue("token"); token != "" {
		completePasswordReset(ctx, token)
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
		ctx.Logger.Error("Failed to fetch user for password reset", "email", emailAddress, "error", err)
		InternalServerError(ctx)
		return
	}
	if user == nil {
		RenderResetPage(ctx, "We could not find any user with that email address.")
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

// completePasswordReset handles the new password submitted from the verification page
func completePasswordReset(ctx *server.Context, token string) {
	verification, err := ctx.State.Verifications.ByToken(token, "User")
	if err != nil {
		ctx.Logger.Error("Failed to fetch verification for password reset", "error", err)
		InternalServerError(ctx)
		return
	}
	if verification == nil {
		NotFound(ctx)
		return
	}
	if verification.Type != constants.VerificationTypePassword {
		NotFound(ctx)
		return
	}

	password := ctx.Request.FormValue("password")
	passwordMatch := ctx.Request.FormValue("password_match")

	if password != passwordMatch {
		RenderVerificationPage(ctx, verification, false, true, "The passwords don't match. Please try again!")
		return
	}
	if len(password) < 8 {
		RenderVerificationPage(ctx, verification, false, true, "Please enter a password with at least 8 characters!")
		return
	}

	hashedPassword, err := authentication.CreatePasswordHash(password)
	if err != nil {
		ctx.Logger.Error("Failed to hash new password", "user_id", verification.UserId, "error", err)
		InternalServerError(ctx)
		return
	}

	err = ctx.State.DatabaseTransaction(func(repos *state.Repositories) error {
		_, err := repos.Users.Update(
			&schemas.User{Id: verification.UserId, Bcrypt: hashedPassword},
			"bcrypt",
		)
		if err != nil {
			return err
		}
		return repos.Verifications.DeleteByToken(verification.Token)
	})
	if err != nil {
		ctx.Logger.Error("Failed to reset password", "user_id", verification.UserId, "error", err)
		InternalServerError(ctx)
		return
	}

	ctx.Logger.Info("User successfully reset their password", "user_id", verification.UserId, "username", verification.Username())
	RenderVerificationPage(ctx, verification, true, false, "")
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
	return ctx.State.Redis.Set(
		ctx.Request.Context(),
		fmt.Sprintf("reset_lock:%d", userId), 1,
		12*time.Hour,
	).Err()
}
