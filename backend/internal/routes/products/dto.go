package products

import (
	"time"

	"luny.dev/cherryauctions/internal/models"
)

type ProductImageDTO struct {
	URL     string `json:"url"`
	AltText string `json:"alt"`
}

type ProfileDTO struct {
	Name  string  `json:"name"`
	Email *string `json:"email"`
}

type ProductDTO struct {
	ID                  uint       `json:"id"`
	Name                string     `json:"name"`
	StartingBid         float64    `json:"starting_bid"`
	StepBidType         string     `json:"step_bid_type"`
	StepBidValue        float64    `json:"step_bid_value"`
	BINPrice            float64    `json:"bin_price"`
	Description         string     `json:"description"`
	ThumbnailURL        string     `json:"thumbnail_url"`
	AllowsUnratedBuyers bool       `json:"allows_unrated_buyers"`
	AutoExtendsTime     bool       `json:"auto_extends_time"`
	CreatedAt           time.Time  `json:"created_at"`
	ExpiredAt           time.Time  `json:"expired_at"`
	Seller              ProfileDTO `json:"seller"`
}

func ToProfileDTO(m models.User) ProfileDTO {
	return ProfileDTO{
		Name:  m.Name,
		Email: m.Email,
	}
}

func ToProductImageDTO(m models.ProductImage) ProductImageDTO {
	return ProductImageDTO{
		URL:     m.URL,
		AltText: m.AltText,
	}
}

func ToProductDTO(m *models.Product) ProductDTO {
	return ProductDTO{
		ID:                  m.ID,
		Name:                m.Name,
		StartingBid:         m.StartingBid,
		StepBidType:         m.StepBidType,
		StepBidValue:        m.StepBidValue,
		BINPrice:            m.BINPrice,
		Description:         m.Description,
		ThumbnailURL:        m.ThumbnailURL,
		AllowsUnratedBuyers: m.AllowsUnratedBuyers,
		AutoExtendsTime:     m.AutoExtendsTime,
		CreatedAt:           m.CreatedAt,
		ExpiredAt:           m.ExpiredAt,
		Seller:              ToProfileDTO(m.Seller),
	}
}

type GetProductsQuery struct {
	Query   string `form:"query" json:"query"`
	Page    int    `form:"page" binding:"number,gt=0,omitempty" json:"page"`
	PerPage int    `form:"per_page" binding:"number,gt=0,omitempty" json:"per_page"`
}

type GetProductsResponse struct {
	Data       []ProductDTO `json:"data"`
	Total      int64        `json:"total"`
	TotalPages int          `json:"total_pages"`
	Page       int          `json:"page"`
	PerPage    int          `json:"per_page"`
}
