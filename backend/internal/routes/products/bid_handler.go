package products

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/internal/logging"
	"luny.dev/cherryauctions/internal/models"
	"luny.dev/cherryauctions/internal/routes/shared"
	"luny.dev/cherryauctions/internal/services"
)

// PostBid godoc
//
//	@summary		Creates a bid on a product
//	@description	Bids on a product, hopefully with some race-condition management.
//	@param			body	body	products.PostBidBody	true	"Bid body"
//	@tags			products
//	@security		ApiKeyAuth
//	@accept			json
//	@produce		json
//	@success		201	{object}	shared.MessageResponse	"Successful bid"
//	@failure		400	{object}	shared.ErrorResponse	"Bad request"
//	@failure		401	{object}	shared.ErrorResponse	"User is unauthenticated"
//	@failure		403	{object}	shared.ErrorResponse	"User can not bid for other reasons"
//	@failure		409	{object}	shared.ErrorResponse	"Race condition, and you lost"
//	@failure		500	{object}	shared.ErrorResponse	"Server could not finish the request"
//	@router			/products/{id}/bids [POST]
func (h *ProductsHandler) PostBid(g *gin.Context) {
	ctx := g.Request.Context()
	claims, _ := g.Get("claims")
	sub := claims.(*services.JWTSubject)

	paramId := g.Param("id")
	id, err := strconv.ParseInt(paramId, 10, 0)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error(), "id": paramId})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "bad request"})
		return
	}

	var body PostBidBody
	if err := g.ShouldBind(&body); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "bad request"})
		return
	}

	lastBid := models.Bid{}
	newBid := models.Bid{}
	product := models.Product{}
	err = h.ProductRepo.CreateBid(ctx, uint(id), sub.UserID, body.BidAmount, &lastBid, &newBid, &product)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusConflict, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusConflict, shared.ErrorResponse{Error: err.Error()})
		return
	}

	// Setup the bid email
	newBid.User = models.User{Name: &sub.Name, Email: &sub.Email}

	response := shared.MessageResponse{Message: "successful bid"}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusCreated, "body": body, "last_bid": lastBid, "response": response})
	h.MailerService.SendBidEmail(&lastBid, &newBid, &product)
	g.JSON(http.StatusCreated, response)
}
