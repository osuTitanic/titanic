package cdn

import "fmt"

type ApiError struct {
	Method     string
	Path       string
	StatusCode int
	Status     string
	ErrorCode  string
	Message    string
	Body       []byte
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func (e *ApiError) Error() string {
	prefix := fmt.Sprintf("cdn: %s %s returned %s", e.Method, e.Path, e.Status)
	if e.ErrorCode != "" && e.Message != "" {
		return fmt.Sprintf("%s: %s: %s", prefix, e.ErrorCode, e.Message)
	}
	if e.Message != "" {
		return fmt.Sprintf("%s: %s", prefix, e.Message)
	}
	return prefix
}
