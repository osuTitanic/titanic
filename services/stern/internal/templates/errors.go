package templates

import (
	"net/http"
)

func InternalServerErrorFallback(response http.ResponseWriter) {
	response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	response.WriteHeader(http.StatusInternalServerError)
	response.Write([]byte("Something went really wrong..."))
	// TODO: Mayhaps add a custom html page as fallback
}
