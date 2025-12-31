package models

import "gorm.io/gorm"

type ChatSession struct {
	gorm.Model
	Product      Product
	ProductID    uint `gorm:"not null"`
	ChatMessages []ChatMessage
}
