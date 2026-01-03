// Package questions provides endpoints for asking questions.
package questions

import (
	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/internal/models"
	"luny.dev/cherryauctions/internal/repositories"
	"luny.dev/cherryauctions/internal/services"
)

type QuestionsHandler struct {
	mailer       *services.MailerService
	middleware   *services.MiddlewareService
	questionRepo *repositories.QuestionRepository
	productRepo  *repositories.ProductRepository
}

func NewQuestionsHandler(
	mailer *services.MailerService,
	middleware *services.MiddlewareService,
	questionRepo *repositories.QuestionRepository,
	productRepo *repositories.ProductRepository,
) *QuestionsHandler {
	return &QuestionsHandler{
		mailer:       mailer,
		middleware:   middleware,
		questionRepo: questionRepo,
		productRepo:  productRepo,
	}
}

func (h *QuestionsHandler) SetupRouter(g *gin.RouterGroup) {
	r := g.Group("/questions")

	r.POST("", h.middleware.AuthorizedRoute(models.ROLE_USER), h.PostQuestion)
	r.PUT("/:id", h.middleware.AuthorizedRoute(models.ROLE_USER), h.PutQuestion)
}
