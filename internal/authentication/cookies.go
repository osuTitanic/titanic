package authentication

import (
	"net/http"
	"strings"
	"time"

	"github.com/osuTitanic/titanic-go/internal/config"
)

const (
	WebsiteSessionCookieName = "titanic_session"
)

func ResolveCookieDomain(cfg *config.Config) string {
	if cfg == nil || cfg.Debug {
		return ""
	}

	if strings.Contains(cfg.DomainName, "localhost") || strings.Contains(cfg.DomainName, ".local") {
		return ""
	}

	return "." + cfg.DomainName
}

func UseSecureCookies(cfg *config.Config, request *http.Request) bool {
	if request != nil && request.TLS != nil {
		return true
	}

	if cfg == nil {
		return false
	}

	return !cfg.GetAllowInsecureCookies()
}

func NewWebsiteSessionCookie(cfg *config.Config, request *http.Request, token string, maxAge time.Duration) *http.Cookie {
	domain := ResolveCookieDomain(cfg)
	return &http.Cookie{
		Name:     WebsiteSessionCookieName,
		Value:    token,
		Path:     "/",
		Domain:   domain,
		HttpOnly: true,
		MaxAge:   int(maxAge / time.Second),
		Secure:   UseSecureCookies(cfg, request),
		SameSite: http.SameSiteLaxMode,
	}
}

func NewExpiredCookie(name string, cfg *config.Config) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Domain:   ResolveCookieDomain(cfg),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   cfg != nil && !cfg.GetAllowInsecureCookies(),
		SameSite: http.SameSiteLaxMode,
	}
}
