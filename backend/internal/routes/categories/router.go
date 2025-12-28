// Package categories provides endpoints for managing categories.
package categories

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"luny.dev/cherryauctions/internal/models"
	"luny.dev/cherryauctions/internal/repositories"
	"luny.dev/cherryauctions/internal/services"
)

type CategoriesHandler struct {
	DB                *gorm.DB
	CategoryRepo      *repositories.CategoryRepository
	MiddlewareService *services.MiddlewareService
}

func (h *CategoriesHandler) SetupRouter(g *gin.RouterGroup) {
	r := g.Group("/categories")

	r.GET("", h.GetCategories)
	r.POST("", h.MiddlewareService.AuthorizedRoute(models.ROLE_ADMIN), h.PostCategories)
}
