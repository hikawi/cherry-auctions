package products

import (
	"mime/multipart"
	"time"
)

type ProductImageDTO struct {
	URL     string `json:"url"`
	AltText string `json:"alt"`
}

type ProfileDTO struct {
	ID            uint    `json:"id"`
	Name          *string `json:"name"`
	Email         *string `json:"email"`
	AvatarURL     *string `json:"avatar_url"`
	AverageRating float64 `json:"average_rating"`
}

type CategoryDTO struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	ParentID *uint  `json:"parent_id"`
}

type BidDTO struct {
	ID        uint       `json:"id"`
	Price     int64      `json:"price"`
	Automated bool       `json:"automated"`
	Bidder    ProfileDTO `json:"bidder"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type DescriptionChangeDTO struct {
	ID        uint      `json:"id"`
	Changes   string    `json:"changes"`
	CreatedAt time.Time `json:"created_at"`
}

type ProductDTO struct {
	ID                  uint                   `json:"id"`
	Name                string                 `json:"name"`
	StartingBid         int64                  `json:"starting_bid"`
	StepBidValue        int64                  `json:"step_bid_value"`
	BINPrice            *int64                 `json:"bin_price"`
	Description         string                 `json:"description"`
	ThumbnailURL        string                 `json:"thumbnail_url"`
	AllowsUnratedBuyers bool                   `json:"allows_unrated_buyers"`
	AutoExtendsTime     bool                   `json:"auto_extends_time"`
	CreatedAt           time.Time              `json:"created_at"`
	ExpiredAt           time.Time              `json:"expired_at"`
	Seller              ProfileDTO             `json:"seller"`
	CurrentHighestBid   *BidDTO                `json:"current_highest_bid"`
	Categories          []CategoryDTO          `json:"categories"`
	DescriptionChanges  []DescriptionChangeDTO `json:"description_changes"`
	DeniedBidders       []ProfileDTO           `json:"denied_bidders"`

	BidsCount  int  `json:"bids_count"`
	IsFavorite bool `json:"is_favorite"`
}

type QuestionDTO struct {
	ID        uint       `json:"id"`
	Content   string     `json:"content"`
	Answer    *string    `json:"answer"`
	User      ProfileDTO `json:"user"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type GetProductsQuery struct {
	Query   string `form:"query" json:"query"`
	Page    int    `form:"page" binding:"number,gt=0,omitempty" json:"page"`
	PerPage int    `form:"per_page" binding:"number,gt=0,omitempty" json:"per_page"`
}

type PostProductIDBody struct {
	ID uint `form:"id" json:"id" binding:"required,number,gt=0"`
}

type GetProductsResponse struct {
	Data       []ProductDTO `json:"data"`
	Total      int64        `json:"total"`
	TotalPages int          `json:"total_pages"`
	Page       int          `json:"page"`
	PerPage    int          `json:"per_page"`
}

type GetTopProductsResponse struct {
	TopBidded   []ProductDTO `json:"top_bids"`
	EndingSoon  []ProductDTO `json:"ending_soon"`
	HighestBids []ProductDTO `json:"highest_bids"`
}

type GetProductDetailsResponse struct {
	ProductDTO
	ProductImages   []ProductImageDTO `json:"product_images"`
	Questions       []QuestionDTO     `json:"questions"`
	Bids            []BidDTO          `json:"bids"`
	SimilarProducts []ProductDTO      `json:"similar_products"`
}

type PostProductBody struct {
	Name          string                  `form:"name" binding:"required,min=2" json:"name"`
	Description   string                  `form:"description" binding:"required,min=50" json:"description"`
	StartingBid   int64                   `form:"starting_bid" binding:"required,number,gt=0" json:"starting_bid"`
	Categories    []uint                  `form:"categories" binding:"required,min=1" json:"categories"`
	ProductImages []*multipart.FileHeader `form:"product_images" binding:"required" json:"product_images"`
	StepBidValue  int64                   `form:"step_bid_value" binding:"required,number,gt=0" json:"step_bid_value"`
	BINPrice      *int64                  `form:"bin_price" binding:"omitempty,number,gt=0" json:"bin_price"`
	AllowsUnrated bool                    `form:"allows_unrated" json:"allows_unrated"`
	AutoExtends   bool                    `form:"auto_extends" json:"auto_extends"`
	ExpiredAt     time.Time               `form:"expired_at" binding:"required,gt" json:"expired_at"`
}

type PostProductDescriptionBody struct {
	Description string `json:"description" form:"description" binding:"required,min=50"`
}

type GetMyProductsQuery struct {
	Page    int `form:"page" binding:"number,gt=0,omitempty" json:"page"`
	PerPage int `form:"per_page" binding:"number,gt=0,omitempty" json:"per_page"`
}

type PostBidBody struct {
	BidAmount int64 `form:"bid" json:"bid" binding:"number,gt=0,required"`
}
