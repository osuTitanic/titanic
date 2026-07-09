package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/osuTitanic/titanic/internal/authentication"
)

func (ctx *Context) IsAuthenticated() bool {
	return ctx != nil && ctx.CurrentUser != nil
}

func (ctx *Context) ResolveAuthentication() {
	cookie, err := ctx.Request.Cookie(authentication.WebsiteSessionCookieName)
	if err != nil {
		return
	}

	session, err := ctx.State.SessionStore.Validate(ctx.Request.Context(), cookie.Value, time.Now())
	if err != nil {
		ctx.Logger.Warn("Failed to validate website session", "error", err)
		ctx.ExpireSessionCookie()
		return
	}
	if session == nil {
		ctx.ExpireSessionCookie()
		return
	}

	user, err := ctx.State.Users.ById(session.UserId)
	if err != nil {
		ctx.Logger.Warn("Failed to load authenticated user", "user_id", session.UserId, "error", err)
		return
	}
	if user == nil {
		ctx.ExpireSessionCookie()
		return
	}

	ctx.CurrentSession = session
	ctx.CurrentUser = user
	ctx.CSRFToken, err = ctx.EnsureCSRFToken()
	if err != nil {
		ctx.Logger.Warn("Failed to load csrf token", "user_id", user.Id, "error", err)
	}
}

func (ctx *Context) EnsureCSRFToken() (string, error) {
	if !ctx.IsAuthenticated() {
		return "", nil
	}

	token, err := ctx.State.CSRFStore.Get(ctx.Request.Context(), ctx.CurrentUser.Id)
	if err != nil {
		return "", err
	}
	if token != "" {
		return token, nil
	}

	return ctx.State.CSRFStore.Upsert(ctx.Request.Context(), ctx.CurrentUser.Id)
}

func (ctx *Context) RefreshCSRFToken() (string, error) {
	if !ctx.IsAuthenticated() {
		return "", nil
	}

	token, err := ctx.State.CSRFStore.Upsert(ctx.Request.Context(), ctx.CurrentUser.Id)
	if err != nil {
		return "", err
	}

	ctx.CSRFToken = token
	return token, nil
}

func (ctx *Context) ValidateCSRF() (bool, error) {
	if !ctx.IsAuthenticated() {
		return false, nil
	}

	token := strings.TrimSpace(ctx.Request.Header.Get("X-CSRF-Token"))
	if token == "" {
		token = strings.TrimSpace(ctx.Request.FormValue("csrf_token"))
	}

	return ctx.State.CSRFStore.Validate(ctx.Request.Context(), ctx.CurrentUser.Id, token)
}

func (ctx *Context) ExpireSessionCookie() {
	if ctx == nil || ctx.Response == nil || ctx.State == nil {
		return
	}

	http.SetCookie(
		ctx.Response,
		authentication.NewExpiredWebsiteSessionCookie(ctx.State.Config, ctx.Request),
	)
}

func (ctx *Context) DeleteCurrentSessionCookie() error {
	if ctx == nil || ctx.State == nil || ctx.Request == nil {
		return nil
	}
	sessionId := ""

	if ctx.CurrentSession != nil {
		sessionId = ctx.CurrentSession.Id
	} else if cookie, err := ctx.Request.Cookie(authentication.WebsiteSessionCookieName); err == nil {
		sessionId = cookie.Value
	}

	if sessionId == "" {
		return nil
	}

	return ctx.State.SessionStore.Delete(ctx.Request.Context(), sessionId)
}

func (ctx *Context) DeleteCurrentCSRFToken() error {
	if !ctx.IsAuthenticated() {
		return nil
	}

	return ctx.State.Redis.Del(
		ctx.Request.Context(),
		authentication.CSRFRedisKey(ctx.CurrentUser.Id),
	).Err()
}
