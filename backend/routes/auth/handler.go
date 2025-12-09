package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/models"
	"luny.dev/cherryauctions/repositories"
	"luny.dev/cherryauctions/routes/shared"
	"luny.dev/cherryauctions/services"
	"luny.dev/cherryauctions/utils"
)

// PostLogin POST /auth/login
//
//	@summary		Logins to an existing account
//	@description	Logins to an account using a username and a password registered with the server.
//	@tags			authentication
//	@accept			json
//	@produce		json
//	@param			credentials	body		auth.LoginRequest	true	"Login credentials"
//	@success		200			{object}	auth.LoginResponse
//	@failure		400			{object}	shared.ErrorResponse
//	@failure		401			{object}	shared.ErrorResponse
//	@failure		404			{object}	shared.ErrorResponse
//	@failure		421			{object}	shared.ErrorResponse
//	@failure		500			{object}	shared.ErrorResponse
//	@router			/auth/login [POST]
func (h *AuthHandler) PostLogin(g *gin.Context) {
	var body LoginRequest
	err := g.ShouldBindBodyWithJSON(&body)
	if err != nil {
		utils.Log(gin.H{"path": g.Request.URL.Path, "status": http.StatusBadRequest, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRepo := repositories.UserRepository{DB: h.DB}

	// Check if it's in the DB yet.
	user, err := userRepo.GetUserByEmail(body.Email)
	if err != nil {
		utils.Log(gin.H{"path": g.Request.URL.Path, "status": http.StatusNotFound, "error": "account doesn't exist"})
		g.AbortWithStatusJSON(http.StatusNotFound, shared.ErrorResponse{Error: "account doesn't exist"})
		return
	}

	// Check against oauth type
	if user.OauthType != "none" || user.Password == nil {
		utils.Log(gin.H{"path": g.Request.URL.Path, "status": http.StatusMisdirectedRequest, "error": "account uses oauth to authenticate"})
		g.AbortWithStatusJSON(http.StatusMisdirectedRequest, shared.ErrorResponse{Error: "account uses oauth to authenticate"})
		return
	}

	// Check the password hash
	ok, err := services.VerifyPassword(*user.Password, body.Password)
	if err != nil {
		utils.Log(gin.H{"path": g.Request.URL.Path, "status": http.StatusInternalServerError, "error": err.Error(), "body": gin.H{"email": body.Email}})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "server can't verify password"})
		return
	}

	if !ok {
		utils.Log(gin.H{"path": g.Request.URL.Path, "status": http.StatusUnauthorized, "error": "wrong password", "body": gin.H{"email": body.Email}})
		g.AbortWithStatusJSON(http.StatusUnauthorized, shared.ErrorResponse{Error: "wrong password"})
		return
	}

	// Generate a JWT key pair.
	accessToken, err := services.SignJWT(user.ID, user.Email)
	if err != nil {
		utils.Log(gin.H{"path": g.Request.URL.Path, "status": http.StatusInternalServerError, "error": err.Error(), "body": gin.H{"email": body.Email}})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "server can't sign jwt"})
		return
	}

	refreshToken, err := utils.GenerateSecretKey(64)
	if err != nil {
		utils.Log(gin.H{"path": g.Request.URL.Path, "status": http.StatusInternalServerError, "error": err.Error(), "body": gin.H{"email": body.Email}})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "server can't generate jwt key pair"})
		return
	}

	// Save the refresh token.
	tokenRepo := repositories.RefreshTokenRepository{DB: h.DB}
	hashedToken := sha256.Sum256(refreshToken)
	_, err = tokenRepo.SaveUserToken(user.ID, base64.URLEncoding.EncodeToString(hashedToken[:]))
	if err != nil {
		utils.Log(gin.H{"path": g.Request.URL.Path, "status": http.StatusInternalServerError, "error": err.Error(), "body": gin.H{"email": body.Email}})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "server can't hash refresh token"})
		return
	}

	cookieSecure, err := strconv.ParseBool(utils.Fatalenv("COOKIE_SECURE"))
	if err != nil {
		utils.Log(gin.H{"path": g.Request.URL.Path, "status": http.StatusInternalServerError, "error": "", "body": gin.H{"email": body.Email}})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "server can't generate jwt key pair"})
		return
	}

	g.SetCookieData(&http.Cookie{
		Name:     "RefreshToken",
		Value:    base64.URLEncoding.EncodeToString(refreshToken),
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24 * 30 * 3),
		Domain:   utils.Fatalenv("DOMAIN"),
		Secure:   cookieSecure,
		SameSite: http.SameSiteNoneMode,
	})
	g.JSON(200, LoginResponse{AccessToken: accessToken})
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
//	@failure		409			{object}	shared.ErrorResponse	"An account with that email already exists"
//	@failure		500			{object}	shared.ErrorResponse	"The request could not be completed due to server faults"
//	@router			/auth/register [POST]
func (h *AuthHandler) PostRegister(g *gin.Context) {
	var body RegisterRequest
	err := g.ShouldBindBodyWithJSON(&body)
	if err != nil {
		utils.Log(gin.H{"path": g.Request.URL.Path, "status": http.StatusBadRequest, "error": err.Error(), "body": body})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: err.Error()})
		return
	}

	userRepo := repositories.UserRepository{DB: h.DB}

	// Check if it's in the DB yet.
	_, err = userRepo.GetUserByEmail(body.Email)
	if err == nil {
		utils.Log(gin.H{"path": g.Request.URL.Path, "status": http.StatusConflict, "error": "account already exists"})
		g.AbortWithStatusJSON(http.StatusConflict, shared.ErrorResponse{Error: "account already exists"})
		return
	}

	// Check password hashes.
	hashedPassword, err := services.HashPassword(body.Password)
	if err != nil {
		utils.Log(gin.H{"path": g.Request.URL.Path, "status": http.StatusInternalServerError, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "server could not hash passowrd"})
		return
	}

	newUser := models.User{
		Name:      body.Name,
		Email:     body.Email,
		Password:  &hashedPassword,
		OauthType: "none",
	}
	err = userRepo.SaveUser(&newUser)
	if err != nil {
		utils.Log(gin.H{"path": g.Request.URL.Path, "status": http.StatusInternalServerError, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "server could not save new account"})
		return
	}

	g.JSON(201, gin.H{"message": "user successfully created"})
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
	// Refresh the access token and rotate the refresh token.
	cookie, err := g.Cookie("RefreshToken")
	if err != nil {
		utils.Log(gin.H{"path": g.Request.URL.Path, "status": http.StatusUnauthorized, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusUnauthorized, shared.ErrorResponse{Error: "refresh token not found"})
		return
	}

	// Check the refresh token.
	tokenRepo := repositories.RefreshTokenRepository{DB: h.DB}
	decodedCookie, err := base64.URLEncoding.DecodeString(cookie)
	if err != nil {
		utils.Log(gin.H{"path": g.Request.URL.Path, "status": http.StatusUnauthorized, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusUnauthorized, shared.ErrorResponse{Error: "invalid refresh token"})
		return
	}

	hashedCookie := sha256.Sum256(decodedCookie)
	savedToken := base64.URLEncoding.EncodeToString(hashedCookie[:])
	token, err := tokenRepo.GetRefreshToken(savedToken)
	if err != nil || token.IsRevoked {
		utils.Log(gin.H{"path": g.Request.URL.Path, "status": http.StatusUnauthorized, "error": "revoked token or non-existent token"})
		g.AbortWithStatusJSON(http.StatusUnauthorized, shared.ErrorResponse{Error: "invalid refresh token"})
		return
	}

	// Make sure the users are checked.
	if token.User.ID == 0 || token.User.Email == "" {
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "this should not be not preloaded"})
		return
	}

	// Invalidate the token.
	_, err = tokenRepo.InvalidateToken(savedToken)
	if err != nil {
		utils.Log(gin.H{"path": g.Request.URL.Path, "status": http.StatusInternalServerError, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "unexpected error while rotating token"})
		return
	}

	// Generate a new JWT key pair.
	accessToken, err := services.SignJWT(token.User.ID, token.User.Email)
	if err != nil {
		utils.Log(gin.H{"path": g.Request.URL.Path, "status": http.StatusInternalServerError, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "server can't sign jwt"})
		return
	}

	refreshToken, err := utils.GenerateSecretKey(64)
	if err != nil {
		utils.Log(gin.H{"path": g.Request.URL.Path, "status": http.StatusInternalServerError, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "server can't generate jwt key pair"})
		return
	}

	// Save the refresh token.
	hashedToken := sha256.Sum256(refreshToken)
	_, err = tokenRepo.SaveUserToken(token.User.ID, base64.URLEncoding.EncodeToString(hashedToken[:]))
	if err != nil {
		utils.Log(gin.H{"path": g.Request.URL.Path, "status": http.StatusInternalServerError, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "server can't hash refresh token"})
		return
	}

	cookieSecure, err := strconv.ParseBool(utils.Fatalenv("COOKIE_SECURE"))
	if err != nil {
		utils.Log(gin.H{"path": g.Request.URL.Path, "status": http.StatusInternalServerError, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "server can't send a cookie"})
		return
	}

	g.SetCookieData(&http.Cookie{
		Name:     "RefreshToken",
		Value:    base64.URLEncoding.EncodeToString(refreshToken),
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24 * 30 * 3),
		Domain:   utils.Fatalenv("DOMAIN"),
		Secure:   cookieSecure,
		SameSite: http.SameSiteNoneMode,
	})
	g.JSON(200, LoginResponse{AccessToken: accessToken})
}

// PostLogout POST /auth/logout
//
//	@summary		Logouts and invalidates the refresh token if available.
//	@description	Logouts, and also invalidates the refresh token. This does not revoke access tokens.
//	@tags			authentication
//	@success		204	{object}	shared.MessageResponse	"Any request, regardless of authentication status"
//	@router			/auth/logout [POST]
func (h *AuthHandler) PostLogout(g *gin.Context) {
	// Fetch the refresh token cookie.
	cookie, err := g.Cookie("RefreshToken")
	if err == nil {
		tokenRepo := repositories.RefreshTokenRepository{DB: h.DB}
		decodedCookie, err := base64.URLEncoding.DecodeString(cookie)
		if err == nil {
			hashedCookie := sha256.Sum256(decodedCookie)
			_, err = tokenRepo.InvalidateToken(base64.URLEncoding.EncodeToString(hashedCookie[:]))

			if err != nil {
				utils.Log(gin.H{"error": "can't invalidate refresh token, ignoring..."})
			}
		}
	}

	cookieSecure, _ := strconv.ParseBool(utils.Fatalenv("COOKIE_SECURE"))
	g.SetCookie("RefreshToken", "", -1, "/", utils.Fatalenv("DOMAIN"), cookieSecure, true)
	g.JSON(204, shared.MessageResponse{Message: "logged out"})
}
