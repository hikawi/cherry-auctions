package models

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name string `gorm:"not null"`

	ParentID       *uint
	ParentCategory *Category  `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	SubCategories  []Category `gorm:"foreignKey:ParentID"`
}
