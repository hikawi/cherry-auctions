package models

import (
	"time"

	"gorm.io/gorm"
)

type SellerSubscription struct {
	gorm.Model
	ExpiredAt time.Time `gorm:"not null;index"`
	UserID    uint      `gorm:"not null"`
	User      User
}
