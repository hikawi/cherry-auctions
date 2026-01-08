package models

import (
	"gorm.io/gorm"
)

type Bid struct {
	gorm.Model
	Price     int64 `gorm:"not null"`
	Automated bool  `gorm:"not null"`
	IsBIN     bool  `gorm:"not null;default:false"`
	ProductID uint
	Product   Product
	UserID    uint
	User      User
}
