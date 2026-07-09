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

	if strings.Contains(cfg.DomainName, "localhost") {
		// Browsers do not allow us to use ".localhost" for some reason
		// so we just don't set the domain at all in this case
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

	if !strings.Contains(request.Host, cfg.DomainName) {
		// If we, for example, access the page through an IP address, we
		// wouldn't be able to authenticate. So this fixes that.
		domain = ""
	}

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

func NewExpiredWebsiteSessionCookie(cfg *config.Config, request *http.Request) *http.Cookie {
	domain := ResolveCookieDomain(cfg)

	if !strings.Contains(request.Host, cfg.DomainName) {
		domain = ""
	}

	return NewExpiredCookie(WebsiteSessionCookieName, domain, cfg)
}

func NewExpiredCookie(name, domain string, cfg *config.Config) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Domain:   domain,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   cfg != nil && !cfg.GetAllowInsecureCookies(),
		SameSite: http.SameSiteLaxMode,
	}
}
