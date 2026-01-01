package products

import (
	"time"

	"luny.dev/cherryauctions/internal/models"
	"luny.dev/cherryauctions/pkg/ranges"
)

type ProductImageDTO struct {
	URL     string `json:"url"`
	AltText string `json:"alt"`
}

type ProfileDTO struct {
	Name  string  `json:"name"`
	Email *string `json:"email"`
}

type CategoryDTO struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	ParentID *uint  `json:"parent_id"`
}

type BidDTO struct {
	ID        uint       `json:"id"`
	Price     float64    `json:"price"`
	Automated bool       `json:"automated"`
	Bidder    ProfileDTO `json:"bidder"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type ProductDTO struct {
	ID                  uint          `json:"id"`
	Name                string        `json:"name"`
	StartingBid         float64       `json:"starting_bid"`
	StepBidType         string        `json:"step_bid_type"`
	StepBidValue        float64       `json:"step_bid_value"`
	BINPrice            float64       `json:"bin_price"`
	Description         string        `json:"description"`
	ThumbnailURL        string        `json:"thumbnail_url"`
	AllowsUnratedBuyers bool          `json:"allows_unrated_buyers"`
	AutoExtendsTime     bool          `json:"auto_extends_time"`
	CreatedAt           time.Time     `json:"created_at"`
	ExpiredAt           time.Time     `json:"expired_at"`
	Seller              ProfileDTO    `json:"seller"`
	CurrentHighestBid   *BidDTO       `json:"current_highest_bid"`
	Categories          []CategoryDTO `json:"categories"`
	BidsCount           int           `json:"bids_count"`
	IsFavorite          bool          `json:"is_favorite"`
}

type QuestionDTO struct {
	ID        uint       `json:"id"`
	Content   string     `json:"content"`
	Answer    *string    `json:"answer"`
	User      ProfileDTO `json:"user"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func ToProfileDTO(m models.User) ProfileDTO {
	return ProfileDTO{
		Name:  m.Name,
		Email: m.Email,
	}
}

func ToCategoryDTO(m models.Category) CategoryDTO {
	return CategoryDTO{
		ID:       m.ID,
		Name:     m.Name,
		ParentID: m.ParentID,
	}
}

func ToQuestionDTO(m models.Question) QuestionDTO {
	var answer *string = nil
	if m.Answer.Valid {
		answer = &m.Answer.String
	}

	return QuestionDTO{
		ID:        m.ID,
		Content:   m.Content,
		Answer:    answer,
		User:      ToProfileDTO(m.User),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func ToQuestionDTOs(questions []models.Question) []QuestionDTO {
	var dtos []QuestionDTO
	for _, question := range questions {
		dtos = append(dtos, ToQuestionDTO(question))
	}
	return dtos
}

func ToBidDTO(m models.Bid) BidDTO {
	return BidDTO{
		ID:        m.ID,
		Price:     m.Price,
		Automated: m.Automated,
		Bidder:    ToProfileDTO(m.User),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func ToBidDTOs(bids []models.Bid) []BidDTO {
	var dtos []BidDTO
	for _, bid := range bids {
		dtos = append(dtos, ToBidDTO(bid))
	}
	return dtos
}

func ToProductImageDTO(m models.ProductImage) ProductImageDTO {
	return ProductImageDTO{
		URL:     m.URL,
		AltText: m.AltText,
	}
}

func ToProductImageDTOs(images []models.ProductImage) []ProductImageDTO {
	var dtos []ProductImageDTO
	for _, img := range images {
		dtos = append(dtos, ToProductImageDTO(img))
	}
	return dtos
}

func ToProductDTO(m *models.Product) ProductDTO {
	var highestBid *BidDTO = nil
	if m.CurrentHighestBid != nil {
		dto := ToBidDTO(*m.CurrentHighestBid)
		highestBid = &dto
	}

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
		CurrentHighestBid:   highestBid,
		BidsCount:           m.BidsCount,
		Categories:          ranges.Each(m.Categories, ToCategoryDTO),
		IsFavorite:          m.IsFavorite,
	}
}

func ToProductDTOs(products []*models.Product) []ProductDTO {
	var dtos []ProductDTO
	for _, product := range products {
		dtos = append(dtos, ToProductDTO(product))
	}
	return dtos
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
