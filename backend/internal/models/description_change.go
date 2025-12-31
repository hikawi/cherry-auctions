package models

import "gorm.io/gorm"

type DescriptionChange struct {
	gorm.Model
	Changes   string `gorm:"not null"`
	ProductID uint
}
