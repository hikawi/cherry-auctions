// Package services provides independent functions that carry out a certain purpose
// that warrant isolation, and not utility enough. These contribute to business logic.
package services

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/utils"
)

type RecaptchaResponse struct {
	Success     bool     `json:"success"`
	Score       float64  `json:"score"`
	Action      string   `json:"action"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
}

func CheckGrecaptcha(token string, g *gin.Context) bool {
	// Run a captcha check.
	form := url.Values{}
	form.Add("secret", utils.Fatalenv("RECAPTCHA_SECRET"))
	form.Add("response", token)
	form.Add("remoteip", g.ClientIP())

	resp, err := http.PostForm(
		"https://www.google.com/recaptcha/api/siteverify",
		form,
	)
	if err != nil {
		utils.Log(gin.H{"path": g.Request.URL.Path, "status": http.StatusConflict, "error": "recaptcha request failed"})
		g.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't verify captcha"})
		return false
	}
	defer resp.Body.Close()

	var result RecaptchaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		g.JSON(500, gin.H{"error": "invalid recaptcha response"})
		return false
	}

	if !result.Success || result.Score < 0.5 || result.Action != "submit" {
		g.JSON(http.StatusForbidden, gin.H{"error": "recaptcha verification failed"})
		return false
	}

	return true
}
