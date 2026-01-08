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

// DeleteBids godoc
//
//	@summary		Clears all bids on a product.
//	@description	Removes all bidders from a product.
//	@tags			products
//	@security		ApiKeyAuth
//	@accept			json
//	@produce		json
//	@success		204	{object}	shared.MessageResponse	"Successful removal"
//	@failure		400	{object}	shared.ErrorResponse	"Invalid ID"
//	@failure		401	{object}	shared.ErrorResponse	"User is unauthenticated"
//	@failure		403	{object}	shared.ErrorResponse	"User is not the seller"
//	@failure		500	{object}	shared.ErrorResponse	"Server could not finish the request"
//	@router			/products/{id}/bids [DELETE]
func (h *ProductsHandler) DeleteBids(g *gin.Context) {
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

	product, err := h.ProductRepo.GetProductByID(ctx, int(id))
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "couldn't find product"})
		return
	}

	if product.SellerID != sub.UserID {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusForbidden, "error": "not the seller"})
		g.AbortWithStatusJSON(http.StatusForbidden, shared.ErrorResponse{Error: "not the seller"})
		return
	}

	_, err = h.ProductRepo.ClearAllBids(ctx, uint(id))
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "couldn't delete bids"})
		return
	}

	response := shared.MessageResponse{Message: "deleted bids"}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusNoContent, "response": response})
	g.JSON(http.StatusNoContent, response)
}

// PostBIN godoc
//
//	@summary		Make a BIN purchase.
//	@description	Make a BIN bid to a product and end it immediately.
//	@tags			products
//	@security		ApiKeyAuth
//	@accept			json
//	@produce		json
//	@success		200	{object}	shared.MessageResponse	"Successful bid"
//	@failure		400	{object}	shared.ErrorResponse	"Invalid ID"
//	@failure		401	{object}	shared.ErrorResponse	"User is unauthenticated"
//	@failure		403	{object}	shared.ErrorResponse	"User is not the seller"
//	@failure		500	{object}	shared.ErrorResponse	"Server could not finish the request"
//	@router			/products/{id}/bin [POST]
func (h *ProductsHandler) PostBIN(g *gin.Context) {
	// ctx := g.Request.Context()
	// claims, _ := g.Get("claims")
	// sub := claims.(*services.JWTSubject)

	// paramId := g.Param("id")
	// id, err := strconv.ParseInt(paramId, 10, 0)
	// if err != nil {
	// 	logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error(), "id": paramId})
	// 	g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "bad request"})
	// 	return
	// }

	// // This is prob inefficient as hell
	// product, err := h.ProductRepo.GetProductByID(ctx, int(id))
	// if err != nil {
	// 	logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
	// 	g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "couldn't find product"})
	// 	return
	// }
}
