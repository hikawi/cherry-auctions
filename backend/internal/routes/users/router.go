package users

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"luny.dev/cherryauctions/internal/models"
	"luny.dev/cherryauctions/internal/repositories"
	"luny.dev/cherryauctions/internal/services"
)

type UsersHandler struct {
	DB                *gorm.DB
	MiddlewareService *services.MiddlewareService
	UserRepo          *repositories.UserRepository
}

func (h *UsersHandler) SetupRouter(r *gin.RouterGroup) {
	g := r.Group("/users")

	g.GET("/me", h.MiddlewareService.AuthorizedRoute(models.ROLE_USER), h.GetMe)
}
