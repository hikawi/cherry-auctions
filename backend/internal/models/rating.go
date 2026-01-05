package models

import "gorm.io/gorm"

type Rating struct {
	gorm.Model
	Rating     float64 `gorm:"not null"`
	ReviewerID uint
	Reviewer   User
	RevieweeID uint
	Reviewee   User
}
