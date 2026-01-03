package questions

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/internal/logging"
	"luny.dev/cherryauctions/internal/routes/shared"
	"luny.dev/cherryauctions/internal/services"
)

func (h *QuestionsHandler) getJWTSubject(g *gin.Context) *services.JWTSubject {
	claimsAny, _ := g.Get("claims")
	claims := claimsAny.(*services.JWTSubject)
	return claims
}

// PostQuestion godoc
//
//	@summary		Posts a question to a product.
//	@description	Creates a new question to a product.
//	@tags			questions
//	@security		ApiKeyAuth
//	@accept			json
//	@produce		json
//	@param			data	body		questions.PostQuestionBody	true	"New question body"
//	@success		201		{object}	shared.MessageResponse		"Added a new question"
//	@failure		400		{object}	shared.ErrorResponse		"When the request is invalid"
//	@failure		401		{object}	shared.ErrorResponse		"When the user is unauthenticated"
//	@failure		403		{object}	shared.ErrorResponse		"When the user is the seller"
//	@failure		500		{object}	shared.ErrorResponse		"The server could not make the request"
//	@router			/questions [POST]
func (h *QuestionsHandler) PostQuestion(g *gin.Context) {
	ctx := g.Request.Context()
	claims := h.getJWTSubject(g)

	var body PostQuestionBody
	if err := g.ShouldBind(&body); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	product, err := h.productRepo.GetProductByID(ctx, int(body.ProductID))
	if err != nil || product.SellerID == claims.UserID {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusForbidden, "error": "can't question a product you sell", "body": body})
		g.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "can't question a product you sell"})
		return
	}

	// Make the question
	err = h.questionRepo.CreateQuestion(ctx, product.ID, claims.UserID, body.Content)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": "can't add a question", "body": body})
		g.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "can't add a question"})
		return
	}

	response := shared.MessageResponse{Message: "created question"}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusCreated, "body": body})
	h.mailer.SendQuestionEmail(&product, body.Content, claims.Name)
	g.JSON(http.StatusCreated, response)
}

// PutQuestion godoc
//
//	@summary		Answers a question.
//	@description	Answers a question to the product.
//	@tags			questions
//	@security		ApiKeyAuth
//	@accept			json
//	@produce		json
//	@param			id		path		int							true	"Question ID"
//	@param			body	body		questions.PutQuestionBody	true	"Question Body"
//	@success		200		{object}	shared.MessageResponse		"Added an answer"
//	@failure		400		{object}	shared.ErrorResponse		"When the request is invalid"
//	@failure		401		{object}	shared.ErrorResponse		"When the user is unauthenticated"
//	@failure		403		{object}	shared.ErrorResponse		"When the user isn't the seller"
//	@failure		500		{object}	shared.ErrorResponse		"The server could not make the request"
//	@router			/questions/{id} [PUT]
func (h *QuestionsHandler) PutQuestion(g *gin.Context) {
	ctx := g.Request.Context()
	claims := h.getJWTSubject(g)
	id, _ := strconv.ParseUint(g.Param("id"), 10, 0)

	var body PutQuestionBody
	if err := g.ShouldBind(&body); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	product, err := h.productRepo.GetProductByID(ctx, int(id))
	if err != nil || product.SellerID != claims.UserID {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusForbidden, "error": "can't answer a product you don't sell", "body": body})
		g.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "can't answer a product you don't sell"})
		return
	}

	rows, err := h.questionRepo.AnswerProductQuestion(ctx, uint(id), body.Answer)
	if rows == 0 || err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": "couldn't save an answer", "body": body})
		g.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "couldn't save an answer"})
		return
	}

	response := shared.MessageResponse{Message: "question answered"}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "body": body, "response": response})
	h.mailer.SendAnswerEmail(uint(id))
	g.JSON(http.StatusOK, response)
}
