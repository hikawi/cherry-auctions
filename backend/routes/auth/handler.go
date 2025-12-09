package auth

import (
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
//	@success		200
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
		utils.Log(gin.H{"path": g.Request.URL.Path, "error": http.StatusBadRequest, "err": err.Error()})
		g.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRepo := repositories.UserRepository{DB: h.DB}

	// Check if it's in the DB yet.
	user, err := userRepo.GetUserByEmail(body.Email)
	if err != nil {
		utils.Log(gin.H{"path": g.Request.URL.Path, "error": http.StatusNotFound, "err": "account doesn't exist"})
		g.AbortWithStatusJSON(http.StatusNotFound, shared.ErrorResponse{Error: "account doesn't exist"})
		return
	}

	// Check against oauth type
	if user.OauthType != "none" || user.Password == nil {
		utils.Log(gin.H{"path": g.Request.URL.Path, "error": http.StatusMisdirectedRequest, "err": "account uses oauth to authenticate"})
		g.AbortWithStatusJSON(http.StatusMisdirectedRequest, shared.ErrorResponse{Error: "account uses oauth to authenticate"})
		return
	}

	// Check the password hash
	ok, err := services.VerifyPassword(*user.Password, body.Password)
	if err != nil {
		utils.Log(gin.H{"path": g.Request.URL.Path, "error": http.StatusInternalServerError, "err": err.Error(), "body": gin.H{"email": body.Email}})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "server can't verify password"})
		return
	}

	if !ok {
		utils.Log(gin.H{"path": g.Request.URL.Path, "error": http.StatusUnauthorized, "err": "wrong password", "body": gin.H{"email": body.Email}})
		g.AbortWithStatusJSON(http.StatusUnauthorized, shared.ErrorResponse{Error: "wrong password"})
		return
	}

	// Generate a JWT key pair.
	accessToken, err := services.SignJWT(user.ID, user.Email)
	refreshToken, err := utils.GenerateSecretKey(64)
	if err != nil {
		utils.Log(gin.H{"path": g.Request.URL.Path, "error": http.StatusInternalServerError, "err": err.Error(), "body": gin.H{"email": body.Email}})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "server can't generate jwt key pair"})
		return
	}

	cookieSecure, err := strconv.ParseBool(utils.Fatalenv("COOKIE_SECURE"))
	g.SetCookieData(&http.Cookie{
		Name:     "RefreshToken",
		Value:    base64.URLEncoding.EncodeToString(refreshToken),
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24 * 30 * 3),
		Domain:   utils.Fatalenv("DOMAIN"),
		Secure:   cookieSecure,
		SameSite: http.SameSiteNoneMode,
	})
	g.SetCookieData(&http.Cookie{
		Name:     "Authorization",
		Value:    accessToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 1),
		Domain:   utils.Fatalenv("DOMAIN"),
		Secure:   cookieSecure,
		SameSite: http.SameSiteNoneMode,
	})
	g.JSON(200, gin.H{})
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
		utils.Log(gin.H{"path": g.Request.URL.Path, "error": http.StatusBadRequest, "err": err.Error(), "body": body})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: err.Error()})
		return
	}

	userRepo := repositories.UserRepository{DB: h.DB}

	// Check if it's in the DB yet.
	_, err = userRepo.GetUserByEmail(body.Email)
	if err == nil {
		utils.Log(gin.H{"path": g.Request.URL.Path, "error": http.StatusConflict, "err": "account already exists"})
		g.AbortWithStatusJSON(http.StatusConflict, shared.ErrorResponse{Error: "account already exists"})
		return
	}

	// Check password hashes.
	hashedPassword, err := services.HashPassword(body.Password)
	if err != nil {
		utils.Log(gin.H{"path": g.Request.URL.Path, "error": http.StatusInternalServerError, "err": err.Error()})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "server could not hash passowrd"})
		return
	}

	newUser := models.User{
		Name:      body.Name,
		Email:     body.Email,
		Password:  &hashedPassword,
		OauthType: "none",
	}
	userRepo.SaveUser(&newUser)
	g.JSON(201, gin.H{"message": "user successfully created"})
}
