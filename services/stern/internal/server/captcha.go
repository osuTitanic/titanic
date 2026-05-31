package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const captchaVerificationUrl = "https://www.google.com/recaptcha/api/siteverify"
const captchaVerificationTimeout = 5 * time.Second

func VerifyRecaptchaResponse(ctx *Context, responseToken string) (bool, error) {
	data := url.Values{
		"secret":   {ctx.State.Config.RecaptchaSecretKey},
		"response": {responseToken},
		"remoteip": {ctx.IP()},
	}

	client := http.Client{Timeout: captchaVerificationTimeout}
	response, err := client.PostForm(captchaVerificationUrl, data)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return false, fmt.Errorf("recaptcha verification returned status %d", response.StatusCode)
	}

	var responseData struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		return false, err
	}
	return responseData.Success, nil
}
