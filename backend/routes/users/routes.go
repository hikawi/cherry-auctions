package users

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"luny.dev/cherryauctions/routes/shared"
)

type UsersHandler struct {
	DB *gorm.DB
}

func (h *UsersHandler) SetupRouter(r *gin.RouterGroup) {
	g := r.Group("/users")
	g.GET("/me", shared.AuthenticatedRoute, h.GetMe)
}
