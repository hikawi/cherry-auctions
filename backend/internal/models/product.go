package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name                string    `gorm:"size:255;not null"`
	StartingBid         float64   `gorm:"type:decimal(10,2);not null"`
	StepBidType         string    `gorm:"check:step_bid_type in ('percentage','value')"`
	StepBidValue        float64   `gorm:"type:decimal(10,2);not null"`
	BINPrice            float64   `gorm:"type:decimal(10,2);not null"`
	Description         string    `gorm:"not null"`
	ThumbnailURL        string    `gorm:"not null"`
	AllowsUnratedBuyers bool      `gorm:"not null;default:true"`
	AutoExtendsTime     bool      `gorm:"not null;default:true"`
	ExpiredAt           time.Time `gorm:"not null"`

	ProductImages []ProductImage
	Categories    []Category `gorm:"many2many:products_categories"`
	Questions     []Question
	SellerID      uint `gorm:"not null"`
	Seller        User

	SearchVector string `gorm:"type:tsvector;index:,type:gin"`
}
