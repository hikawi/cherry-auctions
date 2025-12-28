package auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"luny.dev/cherryauctions/internal/repositories"
	"luny.dev/cherryauctions/internal/services"
)

type AuthHandler struct {
	DB           *gorm.DB
	CookieSecure bool
	Domain       string

	RandomService   *services.RandomService
	JWTService      *services.JWTService
	PasswordService *services.PasswordService
	CaptchaService  *services.CaptchaService

	RefreshTokenRepo *repositories.RefreshTokenRepository
	UserRepo         *repositories.UserRepository
}

func (h *AuthHandler) SetupRouter(group *gin.RouterGroup) {
	router := group.Group("/auth")

	router.POST("/login", h.PostLogin)
	router.POST("/register", h.PostRegister)
	router.POST("/logout", h.PostLogout)
	router.POST("/refresh", h.PostRefresh)
}
