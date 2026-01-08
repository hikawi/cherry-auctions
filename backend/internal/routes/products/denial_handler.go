package products

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/internal/logging"
	"luny.dev/cherryauctions/internal/routes/shared"
	"luny.dev/cherryauctions/internal/services"
)

// PostDenyBidder godoc
//
//	@summary		Deny a bidder.
//	@description	Deny a bidder from a current product.
//	@tags			products
//	@accept			json
//	@produce		json
//	@param			body	body		products.PostDenyBidderBody	true	"Body"
//	@success		200		{object}	shared.MessageResponse		"Successfully denied bidder"
//	@failure		400		{object}	shared.ErrorResponse		"Invalid bidder user ID"
//	@failure		401		{object}	shared.ErrorResponse		"User is unauthorized or not verified"
//	@failure		403		{object}	shared.ErrorResponse		"User is not seller, or denied bidder id is the seller"
//	@failure		404		{object}	shared.ErrorResponse		"Product is not found"
//	@failure		409		{object}	shared.ErrorResponse		"Bidder is already denied"
//	@failure		500		{object}	shared.ErrorResponse		"The server couldn't complete the request"
//	@router			/products/{id}/denials [POST]
func (h *ProductsHandler) PostDenyBidder(g *gin.Context) {
	claimsAny, _ := g.Get("claims")
	claims := claimsAny.(*services.JWTSubject)
	ctx := g.Request.Context()

	paramId := g.Param("id")
	id, err := strconv.ParseUint(paramId, 10, 0)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "invalid product id"})
		return
	}

	product, err := h.ProductRepo.GetProductByID(ctx, int(id))
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusNotFound, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusNotFound, shared.ErrorResponse{Error: "product id not found"})
		return
	}

	if product.SellerID != claims.UserID {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusForbidden, "error": "user is not seller"})
		g.AbortWithStatusJSON(http.StatusForbidden, shared.ErrorResponse{Error: "user is not seller"})
		return
	}

	if product.ExpiredAt.Before(time.Now()) {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": "product is expired"})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "product is expired"})
		return
	}

	var body PostDenyBidderBody
	if err := g.ShouldBind(&body); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "bad request"})
		return
	}

	err = h.ProductRepo.DenyBidder(ctx, product.ID, body.UserID)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "failed to deny bidder in db"})
		return
	}

	response := shared.MessageResponse{Message: "denied bidder"}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "response": response})
	g.JSON(http.StatusOK, response)
}
