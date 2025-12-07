// Package auth adds handlers for routing that involves authentication
// and authorization.
package auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB *gorm.DB
}

// PostLogin POST /auth/login
//
// @summary Logins to an existing account
// @accept json
// @produce json
func (h *AuthHandler) PostLogin(g *gin.Context) {

}
