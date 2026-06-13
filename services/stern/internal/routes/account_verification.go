package routes

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/osuTitanic/titanic-go/internal/constants"
	"github.com/osuTitanic/titanic-go/internal/email"
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/internal/state"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
	"gorm.io/gorm"
)

func RenderVerificationPage(ctx *server.Context, verification *schemas.Verification, success bool, reset bool, errorMessage string) {
	view := templates.VerificationView{
		DefaultView:  buildDefaultView(ctx),
		Verification: verification,
		Success:      success,
		Reset:        reset,
		ErrorMessage: errorMessage,
	}
	ctx.RenderTemplate(http.StatusOK, "pages/account/verification", view)
}

func AccountVerification(ctx *server.Context) {
	if ctx.IsAuthenticated() {
		ctx.Redirect(http.StatusSeeOther, "/")
		return
	}
	query := ctx.Request.URL.Query()

	verificationType, ok := parseVerificationType(query.Get("type"))
	if !ok {
		NotFound(ctx)
		return
	}

	verificationId, ok := parseVerificationId(query.Get("id"))
	if !ok {
		NotFound(ctx)
		return
	}

	verification, err := ctx.State.Verifications.ById(verificationId, "User")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(ctx)
			return
		}
		ctx.Logger.Error("Failed to fetch verification", "verification_id", verificationId, "error", err)
		InternalServerError(ctx)
		return
	}

	token := query.Get("token")
	if token == "" {
		// Let the user know that they have received an email
		RenderVerificationPage(ctx, verification, false, false, "")
		return
	}
	if token != verification.Token {
		ctx.Logger.Warn("Received invalid verification token", "verification_id", verification.Id)
		NotFound(ctx)
		return
	}
	if verificationType != verification.Type {
		NotFound(ctx)
		return
	}
	// Now that we have received a valid token, we can resolve the verification

	switch verification.Type {
	case constants.VerificationTypeActivation:
		if err := activateUserFromVerification(ctx, verification); err != nil {
			ctx.Logger.Error("Failed to activate verification", "verification_id", verification.Id, "user_id", verification.UserId, "error", err)
			InternalServerError(ctx)
			return
		}
		ctx.Logger.Info("Account successfully verified", "user_id", verification.UserId, "username", verification.Username())
		RenderVerificationPage(ctx, verification, true, false, "")

	case constants.VerificationTypePassword:
		// Let the user choose a new password by displaying the reset form
		RenderVerificationPage(ctx, verification, false, true, "")

	default:
		NotFound(ctx)
	}
}

func AccountVerificationResend(ctx *server.Context) {
	if ctx.IsAuthenticated() {
		NotFound(ctx)
		return
	}

	verificationId, ok := parseVerificationId(ctx.Request.URL.Query().Get("id"))
	if !ok {
		NotFound(ctx)
		return
	}

	previousVerification, err := ctx.State.Verifications.ById(verificationId, "User")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(ctx)
			return
		}
		ctx.Logger.Error("Failed to fetch verification for resend", "verification_id", verificationId, "error", err)
		InternalServerError(ctx)
		return
	}

	if previousVerification.IsRecent() {
		RenderVerificationPage(ctx, previousVerification, false, false, "Please wait a few minutes, until you resend the email!")
		return
	}

	newVerification, err := replaceVerification(ctx, previousVerification, previousVerification.Type)
	if err != nil {
		ctx.Logger.Error("Failed to replace verification", "verification_id", previousVerification.Id, "user_id", previousVerification.UserId, "error", err)
		InternalServerError(ctx)
		return
	}
	newVerification.User = previousVerification.User

	if err := sendVerificationEmail(ctx, newVerification); err != nil {
		ctx.Logger.Error("Failed to send verification email", "user_id", newVerification.UserId, "verification_id", newVerification.Id, "error", err)
		RenderVerificationPage(ctx, newVerification, false, false, "Failed to send verification email. Please try again later!")
		return
	}

	ctx.Logger.Info("Resent verification email", "user_id", newVerification.UserId, "username", newVerification.Username())
	ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/account/verification?id=%d", newVerification.Id))
}

// sendVerificationEmail dispatches the correct email for the given verification type
func sendVerificationEmail(ctx *server.Context, verification *schemas.Verification) error {
	switch verification.Type {
	case constants.VerificationTypePassword:
		return sendPasswordResetEmail(ctx, verification)
	default:
		return sendWelcomeEmail(ctx, verification)
	}
}

func activateUserFromVerification(ctx *server.Context, verification *schemas.Verification) error {
	return ctx.State.DatabaseTransaction(func(repos *state.Repositories) error {
		_, err := repos.Users.Update(
			&schemas.User{Id: verification.UserId, Activated: true},
			"activated",
		)
		if err != nil {
			return err
		}
		return repos.Verifications.DeleteByToken(verification.Token)
	})
}

func ensureActivationVerification(ctx *server.Context, user *schemas.User) (*schemas.Verification, bool, error) {
	pending, err := ctx.State.Verifications.ManyByUserIdAndType(user.Id, constants.VerificationTypeActivation)
	if err != nil {
		return nil, false, err
	}

	// Sort pending verifications by sent time, newest first
	sort.Slice(pending, func(i, j int) bool {
		return pending[i].SentAt.After(pending[j].SentAt)
	})

	if len(pending) == 0 {
		// We have no pending verification, so we create a new one
		verification, err := createActivationVerification(ctx, user)
		return verification, err == nil, err
	}

	// Check if the newest verification is still valid,
	// otherwise we replace it with a new one
	verification := pending[0]
	verification.User = user
	if time.Since(verification.SentAt) <= 12*time.Hour {
		return verification, false, nil
	}

	verification, err = replaceVerification(ctx, verification, constants.VerificationTypeActivation)
	if err != nil {
		return nil, false, err
	}

	verification.User = user
	return verification, true, nil
}

func createActivationVerification(ctx *server.Context, user *schemas.User) (*schemas.Verification, error) {
	var verification *schemas.Verification

	// Use transaction to ensure that the verification is created atomically
	err := ctx.State.DatabaseTransaction(func(repos *state.Repositories) error {
		token, err := generateVerificationToken()
		if err != nil {
			return err
		}

		verification, err = repos.Verifications.CreateForUser(
			user.Id,
			constants.VerificationTypeActivation,
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

func replaceVerification(ctx *server.Context, previousVerification *schemas.Verification, verificationType constants.VerificationType) (*schemas.Verification, error) {
	var verification *schemas.Verification

	// Use transaction such that there are no cases where the user has no verification or multiple verifications
	err := ctx.State.DatabaseTransaction(func(repos *state.Repositories) error {
		if err := repos.Verifications.DeleteByToken(previousVerification.Token); err != nil {
			return err
		}

		token, err := generateVerificationToken()
		if err != nil {
			return err
		}

		verification, err = repos.Verifications.CreateForUser(
			previousVerification.UserId,
			verificationType,
			token,
			time.Now(),
		)
		return err
	})
	if err != nil {
		return nil, err
	}

	verification.User = previousVerification.User
	return verification, nil
}

func sendWelcomeEmail(ctx *server.Context, verification *schemas.Verification) error {
	// TODO: maybe move email sender functions to email module
	if verification.User == nil {
		return errors.New("verification email: missing verification user")
	}
	osuBaseUrl := ctx.State.Config.OsuBaseUrl()

	body := fmt.Sprintf(
		constants.EmailWelcomeBody,
		verification.User.Name,
		osuBaseUrl,
		verification.Id,
		verification.Token,
		osuBaseUrl,
	)
	return ctx.State.Email.Send(&email.Message{
		To:       []string{verification.User.Email},
		Subject:  constants.EmailWelcomeSubject,
		TextBody: body,
	})
}

func generateVerificationToken() (string, error) {
	token := make([]byte, 16)
	if _, err := rand.Read(token); err != nil {
		return "", fmt.Errorf("verification token generation: %w", err)
	}
	return hex.EncodeToString(token), nil
}

func parseVerificationType(value string) (constants.VerificationType, bool) {
	// TODO: move this to internal/constants
	if value == "" || value == "activation" {
		return constants.VerificationTypeActivation, true
	}
	if value == "password" {
		return constants.VerificationTypePassword, true
	}
	return 0, false
}

func parseVerificationId(value string) (int, bool) {
	id, err := strconv.Atoi(value)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}
