package models

import "gorm.io/gorm"

type ChatSession struct {
	gorm.Model
	Product      Product
	ProductID    uint `gorm:"not null;uniqueIndex"`
	SellerID     uint `gorm:"not null;index"`
	Seller       User
	BuyerID      uint `gorm:"not null;index"`
	Buyer        User
	ChatMessages []ChatMessage
}
