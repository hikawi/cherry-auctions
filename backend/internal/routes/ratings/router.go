// Package ratings provides ways to set ratings for users
package ratings

import (
	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/internal/models"
	"luny.dev/cherryauctions/internal/repositories"
	"luny.dev/cherryauctions/internal/services"
)

type RatingHandler struct {
	ratingRepo *repositories.RatingRepostory
	middleware *services.MiddlewareService
}

func NewRatingRouter(
	ratingRepo *repositories.RatingRepostory,
	middleware *services.MiddlewareService,
) *RatingHandler {
	return &RatingHandler{
		ratingRepo: ratingRepo,
		middleware: middleware,
	}
}

func (h *RatingHandler) SetupRouter(g *gin.RouterGroup) {
	r := g.Group("/ratings")

	r.POST("", h.middleware.AuthorizedRoute(models.ROLE_USER), h.PostRating)
	r.PUT("/:id", h.middleware.AuthorizedRoute(models.ROLE_USER), h.PutRating)
}
