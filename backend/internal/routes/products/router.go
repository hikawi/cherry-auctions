// Package products provides endpoints for reading and querying products.
package products

import (
	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/internal/models"
	"luny.dev/cherryauctions/internal/repositories"
	"luny.dev/cherryauctions/internal/services"
)

type ProductsHandler struct {
	ProductRepo       *repositories.ProductRepository
	MiddlewareService *services.MiddlewareService
	S3Service         *services.S3Service
	S3PermURL         string
}

func (h *ProductsHandler) SetupRouter(g *gin.RouterGroup) {
	r := g.Group("/products")

	r.GET("", h.MiddlewareService.SoftAuthorizedRoute, h.GetProducts)
	r.GET("/top", h.MiddlewareService.SoftAuthorizedRoute, h.GetProductsTop)
	r.GET("/favorite", h.MiddlewareService.AuthorizedRoute(models.ROLE_USER), h.GetFavoriteProducts)
	r.POST("/favorite", h.MiddlewareService.AuthorizedRoute(models.ROLE_USER), h.PostFavoriteProduct)
	r.POST("", h.MiddlewareService.AuthorizedRoute(models.ROLE_USER), h.PostProduct)
	r.GET("/:id", h.MiddlewareService.SoftAuthorizedRoute, h.GetProductID)
}
