// Package categories provides endpoints for managing categories.
package categories

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoriesHandler struct {
	DB *gorm.DB
}

func (h *CategoriesHandler) SetupRouter(g *gin.RouterGroup) {
	r := g.Group("/categories")

	r.GET("", h.GetCategories)
}
