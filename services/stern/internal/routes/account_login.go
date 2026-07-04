package routes

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/osuTitanic/titanic-go/internal/authentication"
	"github.com/osuTitanic/titanic-go/internal/schemas"
	"github.com/osuTitanic/titanic-go/services/stern/internal/server"
	"github.com/osuTitanic/titanic-go/services/stern/internal/templates"
	"github.com/redis/go-redis/v9"
)

const (
	WebsiteSessionTTL           = 24 * time.Hour
	WebsiteSessionTTLRemembered = 30 * 24 * time.Hour
)

func AccountLoginPage(ctx *server.Context) {
	redirectTarget := sanitizeRedirectTarget(ctx.QueryValue("redirect"))
	if ctx.IsAuthenticated() {
		ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/u/%d", ctx.CurrentUser.Id))
		return
	}
	RenderLoginPage(ctx, "", redirectTarget)
}

func RenderLoginPage(ctx *server.Context, errorMessage string, redirectTarget string) {
	view := templates.LoginView{
		DefaultView:  buildDefaultView(ctx),
		Redirect:     sanitizeRedirectTarget(redirectTarget),
		ErrorMessage: errorMessage,
	}
	ctx.RenderTemplate(http.StatusOK, "pages/account/login", view)
}

func AccountLogin(ctx *server.Context) {
	if err := ctx.Request.ParseForm(); err != nil {
		ctx.Logger.Warn("Failed to parse login form", "error", err)
		RenderLoginPage(ctx, "The specified username or password is incorrect.", "")
		return
	}

	redirect := ctx.Request.FormValue("redirect")
	redirectTarget := sanitizeRedirectTarget(redirect)
	if redirectTarget == "" {
		redirectTarget = "/"
	}

	if ctx.IsAuthenticated() {
		ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/u/%d", ctx.CurrentUser.Id))
		return
	}

	tooManyAttempts, err := hasTooManyLoginAttempts(ctx)
	if err != nil {
		ctx.Logger.Error("Failed to resolve login attempt count", "ip", ctx.IP(), "error", err)
		InternalServerError(ctx)
		return
	}
	if tooManyAttempts {
		RenderLoginPage(ctx, "Too many login attempts. Please wait a minute and try again!", redirectTarget)
		return
	}
	if err := recordLoginAttempt(ctx); err != nil {
		ctx.Logger.Warn("Failed to record login attempt", "ip", ctx.IP(), "error", err)
	}

	identifier := ctx.Request.FormValue("username")
	user, err := resolveLoginUser(ctx, identifier)
	if err != nil {
		ctx.Logger.Error("Failed to resolve login user", "error", err)
		InternalServerError(ctx)
		return
	}

	password := ctx.Request.FormValue("password")
	if user == nil || !authentication.VerifyPasswordHash(password, user.Bcrypt) {
		RenderLoginPage(ctx, "The specified username or password is incorrect.", redirectTarget)
		return
	}

	if !user.Activated {
		// Ensure we have a valid activation & retrieve it
		verification, shouldSendEmail, err := ensureActivationVerification(ctx, user)
		if err != nil {
			ctx.Logger.Error("Failed to prepare account verification", "user_id", user.Id, "error", err)
			InternalServerError(ctx)
			return
		}
		// If we created a new verification, send the email
		if shouldSendEmail {
			if err := sendWelcomeEmail(ctx, verification); err != nil {
				ctx.Logger.Error("Failed to send account verification email", "user_id", user.Id, "verification_id", verification.Id, "error", err)
			}
		}
		ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/account/verification?id=%d", verification.Id))
		return
	}

	remember := ctx.Request.FormValue("remember") != ""
	sessionTTL, persistCookie := resolveWebsiteSessionLifetime(remember)

	session, err := ctx.State.SessionStore.Create(ctx.Request.Context(), user.Id, time.Now(), sessionTTL)
	if err != nil {
		ctx.Logger.Error("Failed to create website session", "user_id", user.Id, "error", err)
		InternalServerError(ctx)
		return
	}

	token, err := ctx.State.CSRFStore.Upsert(ctx.Request.Context(), user.Id)
	if err != nil {
		ctx.State.SessionStore.Delete(ctx.Request.Context(), session.Id)
		ctx.Logger.Error("Failed to create csrf token", "user_id", user.Id, "error", err)
		InternalServerError(ctx)
		return
	}

	cookie := authentication.NewWebsiteSessionCookie(ctx.State.Config, ctx.Request, session.Id, sessionTTL)
	if !persistCookie {
		cookie.MaxAge = 0
		cookie.Expires = time.Time{}
	}
	http.SetCookie(ctx.Response, cookie)

	login := &schemas.Login{
		UserId:  user.Id,
		Ip:      ctx.IP(),
		Version: "web",
	}
	if err := ctx.State.Logins.Create(login); err != nil {
		ctx.Logger.Warn("Failed to record website login", "user_id", user.Id, "error", err)
	}

	ctx.CSRFToken = token
	ctx.Redirect(http.StatusSeeOther, redirectTarget)
}

func resolveWebsiteSessionLifetime(remember bool) (time.Duration, bool) {
	if remember {
		return WebsiteSessionTTLRemembered, true
	}
	return WebsiteSessionTTL, false
}

func resolveLoginUser(ctx *server.Context, identifier string) (*schemas.User, error) {
	identifier = strings.TrimSpace(identifier)
	if identifier == "" {
		return nil, nil
	}

	// Try to resolve by email
	if strings.Contains(identifier, "@") {
		user, err := ctx.State.Users.ByEmail(identifier)
		if err != nil {
			return nil, err
		}
		if user != nil {
			return user, nil
		}
	}

	// Try to resolve by username
	return ctx.State.Users.BySafeName(schemas.ResolveSafeName(identifier))
}

func hasTooManyLoginAttempts(ctx *server.Context) (bool, error) {
	key := "logins:" + ctx.IP()
	attempts, err := ctx.State.Redis.Get(ctx.Request.Context(), key).Int()
	if err == nil {
		return attempts > 30, nil
	}
	if err == redis.Nil {
		return false, nil
	}
	return false, err
}

func recordLoginAttempt(ctx *server.Context) error {
	key := "logins:" + ctx.IP()
	if err := ctx.State.Redis.Incr(ctx.Request.Context(), key).Err(); err != nil {
		return err
	}
	return ctx.State.Redis.Expire(ctx.Request.Context(), key, 30*time.Second).Err()
}
