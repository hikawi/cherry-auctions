package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name                string    `gorm:"size:255;not null"`
	StartingBid         float64   `gorm:"type:decimal(10,2);not null"`
	StepBidType         string    `gorm:"check:step_bid_type in ('percentage','fixed')"`
	StepBidValue        float64   `gorm:"type:decimal(10,2);not null"`
	BINPrice            float64   `gorm:"type:decimal(10,2);not null"`
	Description         string    `gorm:"not null"`
	ThumbnailURL        string    `gorm:"not null"`
	AllowsUnratedBuyers bool      `gorm:"not null;default:true"`
	AutoExtendsTime     bool      `gorm:"not null;default:true"`
	ExpiredAt           time.Time `gorm:"not null"`

	ProductImages       []ProductImage
	Categories          []Category `gorm:"many2many:products_categories"`
	Questions           []Question
	Bids                []Bid
	SellerID            uint `gorm:"not null"`
	Seller              User
	CurrentHighestBid   *Bid
	CurrentHighestBidID *uint `gorm:"default:null"`
	BidsCount           int   `gorm:"default:0;not null"`

	SearchVector string `gorm:"type:tsvector;index:,type:gin"`
}

// Courtesy of AI.
func (p *Product) BeforeSave(tx *gorm.DB) (err error) {
	// We concatenate Name and Description into the SearchVector.
	// to_tsvector transforms the text into searchable tokens.
	// 'simple' dictionary is used to avoid aggressive stemming in multi-language setups.
	content := p.Name + " " + p.Description

	tx.Statement.SetColumn("SearchVector",
		gorm.Expr("to_tsvector('simple', ?)", content),
	)

	return nil
}
