package transactions

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"luny.dev/cherryauctions/internal/logging"
	"luny.dev/cherryauctions/internal/models"
	"luny.dev/cherryauctions/internal/routes/shared"
	"luny.dev/cherryauctions/internal/services"
)

// PostTransaction godoc
//
//	@summary		Creates a transaction.
//	@description	Creates a transaction to link it with a product.
//	@tags			transactions
//	@accept			json
//	@produce		json
//	@security		ApiKeyAuth
//	@param			body	body		transactions.PostTransactionRequest	true	"Transaction data"
//	@success		201		{object}	shared.IDResponse					"Successfully created a transaction"
//	@failure		400		{object}	shared.ErrorResponse				"Invalid body"
//	@failure		401		{object}	shared.ErrorResponse				"Unauthorized"
//	@failure		403		{object}	shared.ErrorResponse				"Only seller can create this transaction"
//	@failure		409		{object}	shared.ErrorResponse				"Already created the transaction for that product"
//	@failure		500		{object}	shared.ErrorResponse				"The server failed to complete the request"
//	@router			/transactions [post]
func (h *TransactionHandler) PostTransaction(g *gin.Context) {
	ctx := g.Request.Context()
	claims, _ := g.Get("claims")
	sub := claims.(*services.JWTSubject)

	var body PostTransactionRequest
	if err := g.ShouldBind(&body); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error(), "body": body})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "bad request"})
		return
	}

	product, err := h.productRepo.GetProductByID(ctx, int(body.ProductID))
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusForbidden, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusForbidden, shared.ErrorResponse{Error: "unknown product"})
		return
	}

	if product.SellerID != sub.UserID {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusForbidden, "error": "can't initiate transaction"})
		g.AbortWithStatusJSON(http.StatusForbidden, shared.ErrorResponse{Error: "can't initiate transaction"})
		return
	}

	if product.ProductState != models.ProductStateEnded || product.CurrentHighestBidID == nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusForbidden, "error": "can't initiate on a product with no winners"})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "can't initiate on a product with no winners"})
		return
	}

	transaction := models.Transaction{
		ProductID:         product.ID,
		BuyerID:           product.CurrentHighestBid.UserID,
		SellerID:          product.SellerID,
		FinalPrice:        product.CurrentHighestBid.Price,
		TransactionStatus: models.TransactionStatusPending,
	}
	err = h.transactionRepo.CreateTransaction(ctx, &transaction)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusForbidden, "error": "already created for the product"})
		g.AbortWithStatusJSON(http.StatusConflict, shared.ErrorResponse{Error: "already created for the product"})
		return
	}

	response := shared.IDResponse{ID: transaction.ID}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusCreated, "response": response})
	g.AbortWithStatusJSON(http.StatusCreated, response)
}

// GetTransactionStatus godoc
//
//	@summary		Gets a transaction status.
//	@description	Gets a transaction status
//	@security		ApiKeyAuth
//	@tags			transactions
//	@accept			json
//	@produce		json
//	@param			id	path		int										true	"Transaction ID"
//	@success		200	{object}	transactions.TransactionStatusResponse	"Successfully queried"
//	@failure		400	{object}	shared.ErrorResponse					"Invalid body"
//	@failure		401	{object}	shared.ErrorResponse					"Unauthorized"
//	@failure		403	{object}	shared.ErrorResponse					"Only seller can create this transaction"
//	@failure		404	{object}	shared.ErrorResponse					"Unknown transaction ID"
//	@failure		500	{object}	shared.ErrorResponse					"The server failed to complete the request"
//	@router			/transactions/{id} [get]
func (h *TransactionHandler) GetTransactionStatus(g *gin.Context) {
	ctx := g.Request.Context()
	id, err := strconv.ParseUint(g.Param("id"), 10, 0)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "invalid id"})
		return
	}

	transaction, err := h.transactionRepo.GetTransactionByID(ctx, uint(id))
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusNotFound, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusNotFound, shared.ErrorResponse{Error: "unknown transaction"})
		return
	}

	response := TransactionStatusResponse{Status: transaction.TransactionStatus}
	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "response": response})
	g.JSON(http.StatusOK, response)
}

// PutTransaction godoc
//
//	@summary		Updates a transaction.
//	@description	Updates a transaction status.
//	@tags			transactions
//	@accept			json
//	@produce		json
//	@security		ApiKeyAuth
//	@param			id		path		int									true	"Transaction ID"
//	@param			body	body		transactions.PutTransactionRequest	true	"Transaction data"
//	@success		201		{object}	shared.IDResponse					"Successfully created a transaction"
//	@failure		400		{object}	shared.ErrorResponse				"Invalid body"
//	@failure		401		{object}	shared.ErrorResponse				"Unauthorized"
//	@failure		403		{object}	shared.ErrorResponse				"Only seller can create this transaction"
//	@failure		404		{object}	shared.ErrorResponse				"Unknown transaction ID"
//	@failure		500		{object}	shared.ErrorResponse				"The server failed to complete the request"
//	@router			/transactions/{id} [put]
func (h *TransactionHandler) PutTransaction(g *gin.Context) {
	ctx := g.Request.Context()
	claims, _ := g.Get("claims")
	sub := claims.(*services.JWTSubject)

	id, err := strconv.ParseUint(g.Param("id"), 10, 0)
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "invalid id"})
		return
	}

	transaction, err := h.transactionRepo.GetTransactionByID(ctx, uint(id))
	if err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusNotFound, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusNotFound, shared.ErrorResponse{Error: "unknown transaction"})
		return
	}

	var body PutTransactionRequest
	if err := g.ShouldBind(&body); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "invalid id"})
		return
	}

	// I want to allow the seller to cancel only if the current status = pending
	// Pending -> Paid -> Delivered -> Completed
	// Only seller can update it to paid/delivered
	// Only buyer can update it to completed (from delivered)
	// If the transaction is already cancelled, nothing to do.
	if transaction.TransactionStatus == models.TransactionStatusCancelled {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusForbidden, "error": "can't update a cancelled transaction"})
		g.AbortWithStatusJSON(http.StatusForbidden, shared.ErrorResponse{Error: "transaction is cancelled"})
		return
	}

	// What's the best way to structure this? Gemini come help me.
	// 1. Identify Roles
	isSeller := sub.UserID == transaction.SellerID
	isBuyer := sub.UserID == transaction.BuyerID

	// 2. State Machine Logic
	switch body.Status {
	case models.TransactionStatusCancelled:
		// Seller can cancel only if Pending
		if isSeller && transaction.TransactionStatus == models.TransactionStatusPending {
			transaction.TransactionStatus = models.TransactionStatusCancelled
		} else {
			logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusForbidden, "error": "can't cancel at this stage"})
			g.AbortWithStatusJSON(http.StatusForbidden, shared.ErrorResponse{Error: "cannot cancel at this stage"})
			return
		}
	case models.TransactionStatusWinnerPaid:
		// Only seller can move from Pending -> Paid
		if isSeller && transaction.TransactionStatus == models.TransactionStatusPending {
			transaction.TransactionStatus = models.TransactionStatusWinnerPaid
		} else {
			logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusForbidden, "error": "only seller can mark as paid from pending"})
			g.AbortWithStatusJSON(http.StatusForbidden, shared.ErrorResponse{Error: "only seller can mark as paid from pending"})
			return
		}

	case models.TransactionStatusDelivered:
		// Only seller can move from Paid -> Delivered
		if isSeller && transaction.TransactionStatus == models.TransactionStatusWinnerPaid {
			transaction.TransactionStatus = models.TransactionStatusDelivered
		} else {
			logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusForbidden, "error": "only seller can mark as delivered after payment"})
			g.AbortWithStatusJSON(http.StatusForbidden, shared.ErrorResponse{Error: "only seller can mark as delivered after payment"})
			return
		}

	case models.TransactionStatusCompleted:
		// Only buyer can move from Delivered -> Completed
		if isBuyer && transaction.TransactionStatus == models.TransactionStatusDelivered {
			transaction.TransactionStatus = models.TransactionStatusCompleted
		} else {
			logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusForbidden, "error": "only buyer can complete the transaction after delivery"})
			g.AbortWithStatusJSON(http.StatusForbidden, shared.ErrorResponse{Error: "only buyer can complete the transaction after delivery"})
			return
		}

	default:
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusBadRequest, "error": "invalid target status"})
		g.AbortWithStatusJSON(http.StatusBadRequest, shared.ErrorResponse{Error: "invalid target status"})
		return
	}

	// 3. Save to DB and Notify Chat
	if _, err := h.transactionRepo.UpdateTransactionStatus(ctx, transaction.ID, transaction.TransactionStatus); err != nil {
		logging.LogMessage(g, logging.LOG_ERROR, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		g.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorResponse{Error: "failed to update"})
		return
	}

	logging.LogMessage(g, logging.LOG_INFO, gin.H{"status": http.StatusOK, "response": shared.IDResponse{ID: uint(id)}})
	h.chatHandler.SendTransactionChangeNotification(transaction.Product.ChatSession.ID, &transaction)
	g.JSON(http.StatusOK, shared.IDResponse{ID: uint(id)})
}
