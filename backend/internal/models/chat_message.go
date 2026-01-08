package models

import "gorm.io/gorm"

type ChatMessage struct {
	gorm.Model
	Sender        User
	SenderID      uint   `gorm:"not null;index"`
	Content       string `gorm:"not null"`
	ChatSessionID uint   `gorm:"not null;index"`
}
