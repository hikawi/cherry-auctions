// Package services provides independent functions that carry out a certain purpose
// that warrant isolation, and not utility enough. These contribute to business logic.
package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"luny.dev/cherryauctions/pkg/closer"
)

const googleRecaptchaSite = "https://www.google.com/recaptcha/api/siteverify"

type RecaptchaResponse struct {
	Success     bool     `json:"success"`
	Score       float64  `json:"score"`
	Action      string   `json:"action"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
}

type CaptchaService struct {
	RecaptchaSecret string
}

var (
	ErrCaptchaCantVerify      error = errors.New("couldn't verify captcha")
	ErrCaptchaInvalidResponse error = errors.New("invalid recaptcha response")
	ErrCaptchaFailed          error = errors.New("recaptcha verification failed")
)

func (s *CaptchaService) CheckGrecaptcha(token string, clientIP string) error {
	// Run a captcha check.
	form := url.Values{}
	form.Add("secret", s.RecaptchaSecret)
	form.Add("response", token)
	form.Add("remoteip", clientIP)

	resp, err := http.PostForm(googleRecaptchaSite, form)
	if err != nil {
		return ErrCaptchaCantVerify
	}
	defer closer.CloseResources(resp.Body)

	var result RecaptchaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ErrCaptchaInvalidResponse
	}

	if !result.Success || result.Score < 0.5 || result.Action != "submit" {
		return ErrCaptchaFailed
	}

	return nil
}
