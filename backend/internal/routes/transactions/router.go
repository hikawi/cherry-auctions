package transactions

import (
	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/internal/models"
	"luny.dev/cherryauctions/internal/repositories"
	"luny.dev/cherryauctions/internal/routes/chat"
	"luny.dev/cherryauctions/internal/services"
)

type TransactionHandler struct {
	transactionRepo   *repositories.TransactionRepository
	productRepo       *repositories.ProductRepository
	ratingRepo        *repositories.RatingRepostory
	middlewareService *services.MiddlewareService
	chatHandler       *chat.ChatHandler
}

func NewTransactionHandler(
	transactionRepo *repositories.TransactionRepository,
	productRepo *repositories.ProductRepository,
	ratingRepo *repositories.RatingRepostory,
	middlewareService *services.MiddlewareService,
	chatHandler *chat.ChatHandler,
) *TransactionHandler {
	return &TransactionHandler{
		transactionRepo:   transactionRepo,
		productRepo:       productRepo,
		ratingRepo:        ratingRepo,
		middlewareService: middlewareService,
		chatHandler:       chatHandler,
	}
}

func (h *TransactionHandler) SetupRouter(g *gin.RouterGroup) {
	r := g.Group("/transactions")

	r.POST("", h.middlewareService.AuthorizedRoute(models.ROLE_USER), h.PostTransaction)
	r.GET("/:id", h.middlewareService.AuthorizedRoute(models.ROLE_USER), h.GetTransactionStatus)
	r.PUT("/:id", h.middlewareService.AuthorizedRoute(models.ROLE_USER), h.PutTransaction)
}
