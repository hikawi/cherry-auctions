package models

import (
	"gorm.io/gorm"
)

type DeniedBidder struct {
	gorm.Model
	ProductID uint
	Product   Product
	User      User
	UserID    uint

	Reason string
}
