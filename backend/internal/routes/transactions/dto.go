package transactions

import "luny.dev/cherryauctions/internal/models"

type PostTransactionRequest struct {
	ProductID uint `json:"product_id" form:"product_id" binding:"gt=0"`
}

type PutTransactionRequest struct {
	Status models.TransactionStatus `json:"status" binding:"required,oneof=pending paid delivered completed cancelled"`
}

type TransactionStatusResponse struct {
	Status models.TransactionStatus `json:"status"`
}
