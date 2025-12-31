package models

import "gorm.io/gorm"

type Bid struct {
	gorm.Model
	Price     float64 `gorm:"not null"`
	Automated bool    `gorm:"not null"`
	ProductID uint
	Product   Product
	UserID    uint
	User      User
}
