package models

import "gorm.io/gorm"

type TransactionStatus string

const (
	TransactionStatusPending    TransactionStatus = "pending"
	TransactionStatusWinnerPaid TransactionStatus = "paid"
	TransactionStatusDelivered  TransactionStatus = "delivered"
	TransactionStatusCompleted  TransactionStatus = "completed"
	TransactionStatusCancelled  TransactionStatus = "cancelled"
)

type Transaction struct {
	gorm.Model
	ProductID         uint `gorm:"uniqueIndex"`
	Product           Product
	BuyerID           uint `gorm:"index"`
	Buyer             User
	SellerID          uint `gorm:"index"`
	Seller            User
	FinalPrice        int64
	TransactionStatus TransactionStatus
}
