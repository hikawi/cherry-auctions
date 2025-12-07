// Package auth adds handlers for routing that involves authentication
// and authorization.
package auth

import (
	"github.com/gin-gonic/gin"
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
func (h *AuthHandler) PostLogin(g *gin.Context) {
}
