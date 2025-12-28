package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/internal/logging"
	"luny.dev/cherryauctions/internal/routes/shared"
)

// PostLogin POST /auth/login
//
//	@summary		Logins to an existing account
//	@description	Logins to an account using a username and a password registered with the server.
//	@tags			authentication
//	@accept			json
//	@produce		json
//	@param			credentials	body		auth.LoginRequest		true	"Login credentials"
//	@success		200			{object}	auth.LoginResponse		"Login successful"
//	@failure		400			{object}	shared.ErrorResponse	"Bad username or password format"
//	@failure		401			{object}	shared.ErrorResponse	"Wrong password"
//	@failure		404			{object}	shared.ErrorResponse	"Account does not exist"
//	@failure		421			{object}	shared.ErrorResponse	"Account uses oauth but tries to login with password"
//	@failure		500			{object}	shared.ErrorResponse	"Server couldn't complete the request"
//	@router			/auth/login [POST]
func (h *AuthHandler) PostLogin(g *gin.Context) {
	ctx := g.Request.Context()

	var body LoginRequest
	err := g.ShouldBindBodyWithJSON(&body)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loggingBody := body
	loggingBody.Password = "[REDACTED]"

	// Check if it's in the DB yet.
	user, err := h.UserRepo.GetUserByEmail(ctx, body.Email)
	if err != nil || user.Email == nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusNotFound, "error": "account doesn't exist", "body": loggingBody})
		g.AbortWithStatusJSON(http.StatusNotFound, shared.ErrorResponse{Error: "account doesn't exist"})
		return
	}

	// Check against oauth type
	if user.OauthType != "none" || user.Password == nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusMisdirectedRequest, "error": "account uses oauth to authenticate", "body": loggingBody})
		g.AbortWithStatusJSON(http.StatusMisdirectedRequest, shared.ErrorResponse{Error: "account uses oauth to authenticate"})
		return
	}

	// Check the password hash
	ok, err := h.PasswordService.VerifyPassword(*user.Password, body.Password)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error(), "body": loggingBody})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "server can't verify password"})
		return
	}

	if !ok {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusUnauthorized, "error": "wrong password", "body": loggingBody})
		g.AbortWithStatusJSON(http.StatusUnauthorized, shared.ErrorResponse{Error: "wrong password"})
		return
	}

	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "message": "user successfully logged in", "body": loggingBody})
	h.assignJWTKeyPair(g, h.DB, loggingBody, user.ID, *user.Email, h.toRoleString(user.Roles))
}

// PostRegister POST /auth/register
//
//	@summary		Registers a new account
//	@description	Registers a new account with the system using a email-password pair.
//	@tags			authentication
//	@accept			json
//	@produce		json
//	@param			credentials	body		auth.RegisterRequest	true	"Register credentials"
//	@success		201			{object}	shared.MessageResponse	"User was successfully registered"
//	@failure		400			{object}	shared.ErrorResponse	"Request body is invalid"
//	@failure		403			{object}	shared.ErrorResponse	"Captcha failed"
//	@failure		409			{object}	shared.ErrorResponse	"An account with that email already exists"
//	@failure		500			{object}	shared.ErrorResponse	"The request could not be completed due to server faults"
//	@router			/auth/register [POST]
func (h *AuthHandler) PostRegister(g *gin.Context) {
	ctx := g.Request.Context()

	var body RegisterRequest
	err := g.ShouldBindBodyWithJSON(&body)
	loggingBody := body
	loggingBody.Password = "[REDACTED]"
	loggingBody.CaptchaToken = "[REDACTED]"

	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error(), "body": loggingBody})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.CaptchaService.CheckGrecaptcha(body.CaptchaToken, g.ClientIP()); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusForbidden, "error": err.Error(), "body": loggingBody})
		return
	}

	// Check if it's in the DB yet.
	_, err = h.UserRepo.GetUserByEmail(ctx, body.Email)
	if err == nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusConflict, "error": "account already exists", "body": loggingBody})
		g.AbortWithStatusJSON(http.StatusConflict, shared.ErrorResponse{Error: "account already exists"})
		return
	}

	// Check password hashes.
	hashedPassword, err := h.PasswordService.HashPassword(body.Password)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error(), "body": loggingBody})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "server could not hash passowrd"})
		return
	}

	_, err = h.UserRepo.RegisterNewUser(ctx, body.Name, body.Email, hashedPassword)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error(), "body": loggingBody})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "server could not save new account"})
		return
	}

	response := shared.MessageResponse{Message: "user successfully registered"}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusCreated, "body": loggingBody, "response": response})
	g.JSON(http.StatusCreated, response)
}

// PostRefresh POST /auth/refresh
//
//	@summary		Refreshs a JWT key pair.
//	@description	Uses the provided refresh token cookie to refresh on another short-lived access token.
//	@tags			authentication
//	@success		204	{object}	shared.MessageResponse	"Any request, regardless of authentication status"
//	@success		200	{object}	auth.LoginResponse		"Refreshed successfully"
//	@failure		401	{object}	shared.ErrorResponse	"Did not attach refresh token"
//	@router			/auth/refresh [POST]
func (h *AuthHandler) PostRefresh(g *gin.Context) {
	ctx := g.Request.Context()

	// Refresh the access token and rotate the refresh token.
	cookie, err := g.Cookie("RefreshToken")
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusUnauthorized, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusUnauthorized, shared.ErrorResponse{Error: "refresh token not found"})
		return
	}

	// Check the refresh token.
	decodedCookie, err := base64.URLEncoding.DecodeString(cookie)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusUnauthorized, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusUnauthorized, shared.ErrorResponse{Error: "invalid refresh token"})
		return
	}

	hashedCookie := sha256.Sum256(decodedCookie)
	savedToken := base64.URLEncoding.EncodeToString(hashedCookie[:])
	token, err := h.RefreshTokenRepo.GetRefreshToken(ctx, savedToken)
	if err != nil || token.IsRevoked {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusUnauthorized, "error": "revoked token or non-existent token"})
		g.AbortWithStatusJSON(http.StatusUnauthorized, shared.ErrorResponse{Error: "invalid refresh token"})
		return
	}

	// Make sure the users are checked.
	if token.User.ID == 0 || token.User.Email == nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": "this should not be preloaded"})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "this should not be not preloaded"})
		return
	}

	// Invalidate the token.
	_, err = h.RefreshTokenRepo.InvalidateToken(ctx, savedToken)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "unexpected error while rotating token"})
		return
	}

	// Generate a new JWT key pair.
	accessToken, err := h.JWTService.SignJWT(token.User.ID, *token.User.Email, h.toRoleString(token.User.Roles))
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "server can't sign jwt"})
		return
	}

	refreshToken, err := h.RandomService.GenerateSecretKey(64)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "server can't generate jwt key pair"})
		return
	}

	// Save the refresh token.
	hashedToken := sha256.Sum256(refreshToken)
	_, err = h.RefreshTokenRepo.SaveUserToken(ctx, token.User.ID, base64.URLEncoding.EncodeToString(hashedToken[:]))
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "server can't hash refresh token"})
		return
	}

	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "message": "returned an access token"})
	g.SetCookieData(&http.Cookie{
		Name:     "RefreshToken",
		Value:    base64.URLEncoding.EncodeToString(refreshToken),
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24 * 30 * 3),
		Domain:   h.Domain,
		Secure:   h.CookieSecure,
		SameSite: http.SameSiteNoneMode,
	})
	g.JSON(http.StatusOK, LoginResponse{AccessToken: accessToken})
}

// PostLogout POST /auth/logout
//
//	@summary		Logouts and invalidates the refresh token if available.
//	@description	Logouts, and also invalidates the refresh token. This does not revoke access tokens.
//	@tags			authentication
//	@success		204	{object}	shared.MessageResponse	"Any request, regardless of authentication status"
//	@router			/auth/logout [POST]
func (h *AuthHandler) PostLogout(g *gin.Context) {
	ctx := g.Request.Context()

	// Fetch the refresh token cookie.
	cookie, err := g.Cookie("RefreshToken")
	if err == nil {
		decodedCookie, err := base64.URLEncoding.DecodeString(cookie)
		if err == nil {
			hashedCookie := sha256.Sum256(decodedCookie)
			_, err = h.RefreshTokenRepo.InvalidateToken(ctx, base64.URLEncoding.EncodeToString(hashedCookie[:]))
			if err != nil {
				logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": "can't invalidate refresh token, was ignored", "token": string(decodedCookie)})
			}
		}
	}

	logging.LogMessage(g, logging.LOG_INFO, gin.H{"message": "invalidated refresh token", "status": http.StatusNoContent})
	g.SetCookie("RefreshToken", "", -1, "/", h.Domain, h.CookieSecure, true)
	g.JSON(http.StatusNoContent, shared.MessageResponse{Message: "logged out"})
}
