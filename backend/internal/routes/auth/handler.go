package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/internal/logging"
	"luny.dev/cherryauctions/internal/models"
	"luny.dev/cherryauctions/internal/routes/shared"
	"luny.dev/cherryauctions/internal/services"
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
//	@failure		403			{object}	shared.ErrorResponse	"Account is not verified"
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

	var subscription *time.Time
	if len(user.Subscriptions) > 0 {
		subscription = &user.Subscriptions[0].ExpiredAt
	}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "message": "user successfully logged in", "body": loggingBody})
	h.assignJWTKeyPair(g, loggingBody, user.ID, *user.Name, *user.Email, h.toRoleString(user.Roles), subscription, user.Verified)
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
	fmt.Println(cookie)
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
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": "this should be preloaded"})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "this should be not preloaded"})
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
	var subscription *time.Time
	if len(token.User.Subscriptions) > 0 {
		subscription = &token.User.Subscriptions[0].ExpiredAt
	}
	accessToken, err := h.JWTService.SignJWT(token.User.ID, *token.User.Name, *token.User.Email, h.toRoleString(token.User.Roles), subscription, token.User.Verified)
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
		Expires:  time.Now().Add(time.Hour * 24 * 30 * 7),
		Domain:   h.Domain,
		Secure:   h.CookieSecure,
		HttpOnly: true,
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

// PostVerifyCheck godoc
//
//	@summary		Verifies an OTP code.
//	@description	Verifies a user's OTP code.
//	@tags			authentication
//	@success		200	{object}	shared.MessageResponse	"Verification successfully"
//	@failure		400	{object}	shared.ErrorResponse	"Failed to verify"
//	@failure		401	{object}	shared.ErrorResponse	"Not logged in"
//	@failure		422	{object}	shared.ErrorResponse	"Token is valid but user does not exist"
//	@failure		500	{object}	shared.ErrorResponse	"Internal server error"
//	@router			/auth/verify/check [POST]
func (h *AuthHandler) PostVerifyCheck(g *gin.Context) {
	ctx := g.Request.Context()
	claimsAny, ok := g.Get("claims")
	if !ok {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusUnauthorized, "error": "missing bearer token"})
		g.AbortWithStatusJSON(http.StatusUnauthorized, shared.ErrorResponse{Error: "missing Bearer token"})
		return
	}

	claims := claimsAny.(*services.JWTSubject)
	if claims.Verified {
		logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "message": "already verified", "user_id": claims.UserID})
		g.AbortWithStatusJSON(http.StatusOK, shared.MessageResponse{Message: "already verified"})
		return
	}

	var body PostOTPVerifyBody
	if err := g.ShouldBind(&body); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error(), "user_id": claims.UserID})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "bad request"})
		return
	}

	err := h.OTPService.VerifyOTP(ctx, claims.UserID, body.Code)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error(), "user_id": claims.UserID})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "failed to verify otp"})
		return
	}

	// Update the verified status
	_, err = h.UserRepo.UpdateUserVerified(ctx, claims.UserID, true)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error(), "user_id": claims.UserID})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "failed to update verified status"})
		return
	}

	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "message": "verified successfully"})
	g.JSON(http.StatusOK, shared.MessageResponse{Message: "verified successfully, please refresh for a different access token"})
}

// PostVerify godoc
//
//	@summary		Requests an OTP code for verification.
//	@description	Requests the server to send an OTP code for verification.
//	@tags			authentication
//	@success		200	{object}	shared.MessageResponse	"An OTP code has been sent"
//	@success		204	{object}	shared.MessageResponse	"Already verified"
//	@failure		401	{object}	shared.MessageResponse	"User is unauthenticated"
//	@failure		500	{object}	shared.MessageResponse	"Internal server error"
//	@router			/auth/verify [POST]
func (h *AuthHandler) PostVerify(g *gin.Context) {
	ctx := g.Request.Context()
	claimsAny, ok := g.Get("claims")
	if !ok {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusUnauthorized, "error": "missing bearer token"})
		g.AbortWithStatusJSON(http.StatusUnauthorized, shared.ErrorResponse{Error: "missing Bearer token"})
		return
	}

	claims := claimsAny.(*services.JWTSubject)
	if claims.Verified {
		logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusNoContent, "message": "already verified", "user_id": claims.UserID})
		g.AbortWithStatusJSON(http.StatusNoContent, shared.MessageResponse{Message: "already verified"})
		return
	}

	user := &models.User{
		ID:    claims.UserID,
		Name:  &claims.Name,
		Email: &claims.Email,
	}
	err := h.OTPService.SendOTP(ctx, user)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"error": err.Error(), "user_id": user.ID, "status": http.StatusInternalServerError})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "unable to send otp"})
		return
	}

	response := shared.MessageResponse{Message: "otp sent"}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "user_id": user.ID, "response": response})
	g.JSON(http.StatusOK, response)
}
