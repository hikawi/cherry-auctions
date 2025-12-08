package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/models"
	"luny.dev/cherryauctions/repositories"
	"luny.dev/cherryauctions/routes/shared"
	"luny.dev/cherryauctions/services"
	"luny.dev/cherryauctions/utils"
)

// PostLogin POST /auth/login
//
// @summary Logins to an existing account
// @description Logins to an account using a username and a password registered with the server.
// @tags authentication
// @accept json
// @produce json
// @param credentials body auth.LoginRequest true "Login credentials"
// @success 200
// @failure 400
// @failure 401
// @failure 421
// @router /auth/login [POST]
// func (h *AuthHandler) PostLogin(g *gin.Context) {
// 	var body LoginRequest
// 	err := g.ShouldBindBodyWithJSON(&body)
// 	if err != nil {
// 		utils.Log(gin.H{"path": g.Request.URL.Path, "error": http.StatusBadRequest, "err": err.Error()})
// 		g.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
// 		return
// 	}
// }

// type LoginRequest struct {
// 	Email    string `json:"email" binding:"required,email"`
// 	Password string `json:"password" binding:"required,min=8,max=64"`
// }

// PostRegister POST /auth/register
//
// @summary Registers a new account
// @description Registers a new account with the system using a email-password pair.
// @tags authentication
// @accept json
// @produce json
// @param credentials body auth.RegisterRequest true "Register credentials"
// @success 201
// @failure 400
// @failure 401
// @failure 422
// @failure 409
// @router /auth/register [POST]
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
	g.JSON(201, newUser)
}
