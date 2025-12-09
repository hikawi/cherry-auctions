package auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB *gorm.DB
}

func (h *AuthHandler) SetupRouter(group *gin.RouterGroup) {
	router := group.Group("/auth")

	router.POST("/login", h.PostLogin)
	router.POST("/register", h.PostRegister)
	router.POST("/logout", h.PostLogout)
	router.POST("/refresh", h.PostRefresh)
}
