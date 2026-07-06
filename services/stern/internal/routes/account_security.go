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
	"github.com/osuTitanic/titanic-go/internal/location"
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/internal/state"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
)

func AccountSecurity(ctx *server.Context) {
	if !ctx.RequireLogin() {
		return
	}
	renderSecurityPage(ctx, "", "")
}

func AccountSecurityUpdate(ctx *server.Context) {
	if !ctx.RequireLogin() {
		return
	}
	if err := ctx.Request.ParseForm(); err != nil {
		renderSecurityPage(ctx, "", "Please enter your current password!")
		return
	}

	currentPassword := ctx.FormValue("current-password")
	if currentPassword == "" {
		renderSecurityPage(ctx, "", "Please enter your current password!")
		return
	}
	if !authentication.VerifyPasswordHash(currentPassword, ctx.CurrentUser.Bcrypt) {
		renderSecurityPage(ctx, "", "Your password was incorrect. Please try again!")
		return
	}

	newEmail := ctx.FormValue("new-email")
	emailConfirm := ctx.FormValue("email-confirm")
	if newEmail != "" && emailConfirm != "" {
		changeEmailAddress(ctx, newEmail, emailConfirm)
		return
	}

	newPassword := ctx.FormValue("new-password")
	passwordConfirm := ctx.FormValue("password-confirm")
	if newPassword != "" && passwordConfirm != "" {
		changePassword(ctx, newPassword, passwordConfirm)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/account/profile")
}

func changeEmailAddress(ctx *server.Context, newEmail, emailConfirm string) {
	newEmail = strings.ToLower(strings.TrimSpace(newEmail))
	emailConfirm = strings.ToLower(strings.TrimSpace(emailConfirm))

	if newEmail != emailConfirm {
		renderSecurityPage(ctx, "", "The emails don't match. Please try again!")
		return
	}
	if newEmail == ctx.CurrentUser.Email {
		renderSecurityPage(ctx, "", "You're already using that email!")
		return
	}

	existing, err := ctx.State.Users.ByEmail(newEmail)
	if err != nil {
		ctx.Logger.Error("Failed to check email availability", "error", err)
		InternalServerError(ctx)
		return
	}
	if existing != nil {
		renderSecurityPage(ctx, "", "There already is a user with that email. Please choose another one, or reset your password!")
		return
	}

	// Notify the current (old) email address before it is replaced
	if err := sendEmailAddressChanged(ctx, ctx.CurrentUser); err != nil {
		ctx.Logger.Warn("Failed to send email-changed notification", "user", ctx.CurrentUser.Id, "error", err)
	}

	// Deactivate the account, change the email & create a reactivation token
	var verification *schemas.Verification

	// Use transaction to ensure that the email is only changed if the verification token is successfully created
	// Otherwise we'd end up in a state where the email is changed but the user cannot reactivate their account
	err = ctx.State.DatabaseTransaction(func(repos *state.Repositories) error {
		_, err := repos.Users.Update(
			&schemas.User{Id: ctx.CurrentUser.Id, Activated: false, Email: newEmail},
			"activated", "email",
		)
		if err != nil {
			return err
		}

		token, err := generateVerificationToken()
		if err != nil {
			return err
		}

		verification, err = repos.Verifications.CreateForUser(
			ctx.CurrentUser.Id,
			constants.VerificationTypeActivation,
			token,
			time.Now(),
		)
		return err
	})
	if err != nil {
		ctx.Logger.Error("Failed to change email address", "user", ctx.CurrentUser.Id, "error", err)
		InternalServerError(ctx)
		return
	}

	ctx.CurrentUser.Email = newEmail
	ctx.CurrentUser.Activated = false
	verification.User = ctx.CurrentUser

	if err := sendReactivateEmail(ctx, verification); err != nil {
		ctx.Logger.Error("Failed to send reactivation email", "user", ctx.CurrentUser.Id, "error", err)
	}

	// The account is now deactivated -> the current session should be logged out
	logoutCurrentSession(ctx)
	ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/account/verification?id=%d", verification.Id))
}

func changePassword(ctx *server.Context, newPassword, passwordConfirm string) {
	if newPassword != passwordConfirm {
		renderSecurityPage(ctx, "", "The passwords don't match. Please try again!")
		return
	}

	hashedPassword, err := authentication.CreatePasswordHash(newPassword)
	if err != nil {
		ctx.Logger.Error("Failed to hash new password", "user", ctx.CurrentUser.Id, "error", err)
		InternalServerError(ctx)
		return
	}

	_, err = ctx.State.Users.Update(
		&schemas.User{Id: ctx.CurrentUser.Id, Bcrypt: hashedPassword},
		"bcrypt",
	)
	if err != nil {
		ctx.Logger.Error("Failed to update password", "user", ctx.CurrentUser.Id, "error", err)
		InternalServerError(ctx)
		return
	}
	ctx.CurrentUser.Bcrypt = hashedPassword

	if err := sendPasswordChangedEmail(ctx, ctx.CurrentUser); err != nil {
		ctx.Logger.Warn("Failed to send password-changed notification", "user", ctx.CurrentUser.Id, "error", err)
	}
	renderSecurityPage(ctx, "Your password was updated.", "")
}

func renderSecurityPage(ctx *server.Context, infoMessage, errorMessage string) {
	logins, err := ctx.State.Logins.FetchMany(ctx.CurrentUser.Id, 5, 0)
	if err != nil {
		ctx.Logger.Error("Failed to fetch logins", "user", ctx.CurrentUser.Id, "error", err)
		logins = nil
	}

	view := templates.SettingsSecurityView{
		DefaultView:  buildDefaultView(ctx),
		InfoMessage:  infoMessage,
		ErrorMessage: errorMessage,
		IrcUsername:  strings.ReplaceAll(ctx.CurrentUser.Name, " ", "_"),
		Logins:       buildSecurityLogins(ctx, logins),
	}
	ctx.RenderTemplate(http.StatusOK, "pages/account/settings_security", view)
}

func buildSecurityLogins(ctx *server.Context, logins []*schemas.Login) []*templates.SecurityLogin {
	entries := make([]*templates.SecurityLogin, 0, len(logins))
	for _, login := range logins {
		// Logins from the website chat originate from a local IP over IRC
		isWebIrc := location.IsLocalIP(login.Ip) && strings.EqualFold(login.Version, "irc")
		country := "Unknown"

		if !isWebIrc {
			resolved, err := ctx.State.Location.Resolve(login.Ip)
			if err != nil {
				ctx.Logger.Warn("Failed to resolve login location", "ip", login.Ip, "error", err)
			}
			// Resolve always returns a location no matter what, so this is fine
			country = resolved.CountryName
		}

		entries = append(entries, &templates.SecurityLogin{
			Time:     login.Time,
			IsWebIrc: isWebIrc,
			Country:  country,
			Ip:       login.Ip,
			Version:  login.Version,
		})
	}
	return entries
}

func logoutCurrentSession(ctx *server.Context) {
	if err := ctx.DeleteCurrentSessionCookie(); err != nil {
		ctx.Logger.Warn("Failed to delete website session", "user", ctx.CurrentUser.Id, "error", err)
	}
	if err := ctx.DeleteCurrentCSRFToken(); err != nil {
		ctx.Logger.Warn("Failed to delete csrf token", "user", ctx.CurrentUser.Id, "error", err)
	}
	ctx.ExpireSessionCookie()
}

func sendReactivateEmail(ctx *server.Context, verification *schemas.Verification) error {
	if verification.User == nil {
		return errors.New("reactivate email: missing verification user")
	}
	osuBaseUrl := ctx.State.Config.OsuBaseUrl()
	body := fmt.Sprintf(
		constants.EmailReactivateBody,
		verification.User.Name,
		osuBaseUrl,
		verification.Id,
		verification.Token,
		osuBaseUrl,
	)
	return ctx.State.Email.Send(&email.Message{
		To:       []string{verification.User.Email},
		Subject:  constants.EmailReactivateSubject,
		TextBody: body,
	})
}

func sendPasswordChangedEmail(ctx *server.Context, user *schemas.User) error {
	osuBaseUrl := ctx.State.Config.OsuBaseUrl()
	body := fmt.Sprintf(constants.EmailPasswordChangedBody, user.Name, osuBaseUrl, osuBaseUrl)
	return ctx.State.Email.Send(&email.Message{
		To:       []string{user.Email},
		Subject:  constants.EmailPasswordChangedSubject,
		TextBody: body,
	})
}

func sendEmailAddressChanged(ctx *server.Context, user *schemas.User) error {
	osuBaseUrl := ctx.State.Config.OsuBaseUrl()
	body := fmt.Sprintf(constants.EmailAddressChangedBody, user.Name, osuBaseUrl)
	return ctx.State.Email.Send(&email.Message{
		To:       []string{user.Email},
		Subject:  constants.EmailAddressChangedSubject,
		TextBody: body,
	})
}
