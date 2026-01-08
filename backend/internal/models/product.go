package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductState string

const (
	ProductStateActive    ProductState = "active"
	ProductStateEnded     ProductState = "ended"
	ProductStateExpired   ProductState = "expired"
	ProductStateCancelled ProductState = "cancelled"
)

type Product struct {
	gorm.Model
	Name                string       `gorm:"size:255;not null"`
	StartingBid         int64        `gorm:"type:bigint;not null"`
	StepBidValue        int64        `gorm:"type:bigint;not null"`
	BINPrice            *int64       `gorm:"type:bigint"`
	Description         string       `gorm:"not null"`
	ThumbnailURL        string       `gorm:"not null"`
	AllowsUnratedBuyers bool         `gorm:"not null;default:true"`
	AutoExtendsTime     bool         `gorm:"not null;default:true"`
	ExpiredAt           time.Time    `gorm:"not null"`
	EmailSent           bool         `gorm:"not null;default:false"`
	ProductState        ProductState `gorm:"not null;default:active"`

	ProductImages      []ProductImage `gorm:"foreignKey:ProductID"`
	Categories         []Category     `gorm:"many2many:products_categories"`
	Questions          []Question
	Bids               []Bid
	DeniedBidders      []DeniedBidder
	DescriptionChanges []DescriptionChange `gorm:"foreignKey:ProductID"`

	SellerID uint `gorm:"not null"`
	Seller   User

	CurrentHighestBid   *Bid
	CurrentHighestBidID *uint `gorm:"default:null"`

	BidsCount    int    `gorm:"default:0;not null"`
	SearchVector string `gorm:"type:tsvector;index:,type:gin"`
	IsFavorite   bool   `gorm:"-"`
}

// Courtesy of AI.
func (p *Product) BeforeSave(tx *gorm.DB) (err error) {
	content := p.Name + " " + p.Description

	tx.Statement.SetColumn(
		"SearchVector",
		gorm.Expr("to_tsvector('simple', ?)", content),
	)

	return nil
}
