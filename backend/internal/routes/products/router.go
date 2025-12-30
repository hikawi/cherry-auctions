// Package products provides endpoints for reading and querying products.
package products

import (
	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/internal/repositories"
)

type ProductsHandler struct {
	ProductRepo *repositories.ProductRepository
}

func (h *ProductsHandler) SetupRouter(g *gin.RouterGroup) {
	r := g.Group("/products")

	r.GET("", h.GetProducts)
}
