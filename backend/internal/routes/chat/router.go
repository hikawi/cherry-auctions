// Package chat provides endpoints for managing chat sessions between 2 users
package chat

import (
	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/internal/models"
	"luny.dev/cherryauctions/internal/repositories"
	"luny.dev/cherryauctions/internal/services"
)

type ChatHandler struct {
	middlewareService *services.MiddlewareService
	s3Service         *services.S3Service
	chatSessionRepo   *repositories.ChatSessionRepository
	productRepo       *repositories.ProductRepository
	s3PermURL         string
}

func NewChatHandler(
	middlewareService *services.MiddlewareService,
	s3Service *services.S3Service,
	chatSessionRepo *repositories.ChatSessionRepository,
	productRepo *repositories.ProductRepository,
	s3PermURL string,
) *ChatHandler {
	return &ChatHandler{
		middlewareService: middlewareService,
		s3Service:         s3Service,
		chatSessionRepo:   chatSessionRepo,
		productRepo:       productRepo,
		s3PermURL:         s3PermURL,
	}
}

func (h *ChatHandler) SetupRouter(g *gin.RouterGroup) {
	r := g.Group("/chat")

	r.GET("", h.middlewareService.AuthorizedRoute(models.ROLE_USER), h.GetChatSessions)
	r.POST("", h.middlewareService.AuthorizedRoute(models.ROLE_USER), h.CreateChatSession)
	r.GET("/:id", h.middlewareService.AuthorizedRoute(models.ROLE_USER), h.GetChatMessages)
	r.POST("/:id", h.middlewareService.AuthorizedRoute(models.ROLE_USER), h.PostChatMessage)
	r.GET("/stream", h.middlewareService.InjectAuthQuery, h.middlewareService.AuthorizedRoute(models.ROLE_USER), h.GetChatStream)
}
