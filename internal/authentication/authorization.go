package authentication

import (
	"encoding/base64"
	"errors"
	"strings"
)

var ErrInvalidBasicCredentials = errors.New("authentication: invalid basic credentials")

type Authorization struct {
	Scheme string
	Data   string
}

func ParseAuthorization(header string) Authorization {
	if header == "" || !strings.Contains(header, " ") {
		return Authorization{}
	}

	scheme, data, _ := strings.Cut(header, " ")
	return Authorization{
		Scheme: strings.ToLower(strings.TrimSpace(scheme)),
		Data:   strings.TrimSpace(data),
	}
}

func ParseBasicCredentials(encoded string) (username string, password string, err error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", ErrInvalidBasicCredentials
	}

	username, password, ok := strings.Cut(string(decoded), ":")
	if !ok {
		return "", "", ErrInvalidBasicCredentials
	}

	return username, password, nil
}
