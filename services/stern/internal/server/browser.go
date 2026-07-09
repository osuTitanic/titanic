package server

import (
	"regexp"
	"strings"
)

var (
	chromePattern  = regexp.MustCompile(`Chrome/([7-9][0-9]|[1-9][0-9]{2,})`)
	firefoxPattern = regexp.MustCompile(`Firefox/([6-9][0-9]|[1-9][0-9]{2,})`)
	safariPattern  = regexp.MustCompile(`Version/(1[2-9]|[2-9][0-9])\.\d+`)
	edgePattern    = regexp.MustCompile(`Edg/(79|[8-9][0-9]|[1-9][0-9]{2,})`)
	operaPattern   = regexp.MustCompile(`OPR/([6-9][0-9]|[1-9][0-9]{2,})`)
	edgeOrOpera    = regexp.MustCompile(`Edg/|OPR/`)
)

// IsModernBrowser checks if the given user agent string belongs to a modern browser
func IsModernBrowser(userAgent string) bool {
	switch {
	case chromePattern.MatchString(userAgent) && !edgeOrOpera.MatchString(userAgent):
		// Chrome 70+ (excluding Edge and Opera)
		return true
	case firefoxPattern.MatchString(userAgent):
		// Firefox 60+
		return true
	case safariPattern.MatchString(userAgent) && strings.Contains(userAgent, "Safari") && !strings.Contains(userAgent, "Chrome"):
		// Safari 12+ (excluding Chrome)
		return true
	case edgePattern.MatchString(userAgent):
		// Edge (Chromium) 79+
		return true
	case operaPattern.MatchString(userAgent):
		// Opera 60+
		return true
	default:
		return false
	}
}

// IsInternetExplorer checks if the given user agent string belongs to Internet Explorer
func IsInternetExplorer(userAgent string) bool {
	userAgent = strings.ToLower(userAgent)
	return strings.Contains(userAgent, "msie") || strings.Contains(userAgent, "trident/")
}
