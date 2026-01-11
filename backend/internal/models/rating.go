package models

import "gorm.io/gorm"

type Rating struct {
	gorm.Model
	Rating     uint   `gorm:"not null"`
	Feedback   string `gorm:"not null"`
	ProductID  uint   `gorm:"not null;index"`
	Product    Product
	ReviewerID uint `gorm:"not null;index"`
	Reviewer   User
	RevieweeID uint `gorm:"not null;index"`
	Reviewee   User
}
