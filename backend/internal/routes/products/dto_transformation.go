package products

import (
	"luny.dev/cherryauctions/internal/models"
	"luny.dev/cherryauctions/pkg/ranges"
)

func ToProfileDTO(m models.User) ProfileDTO {
	return ProfileDTO{
		ID:            m.ID,
		Name:          m.Name,
		Email:         m.Email,
		AvatarURL:     m.AvatarURL,
		AverageRating: m.AverageRating,
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
		ProductState:        string(m.ProductState),
		DeniedBidders: ranges.Each(m.DeniedBidders, func(bidder models.DeniedBidder) ProfileDTO {
			return ToProfileDTO(bidder.User)
		}),
		DescriptionChanges: ranges.Each(m.DescriptionChanges, func(m models.DescriptionChange) DescriptionChangeDTO {
			return DescriptionChangeDTO{
				ID:        m.ID,
				Changes:   m.Changes,
				CreatedAt: m.CreatedAt,
			}
		}),
		IsFavorite: m.IsFavorite,
	}
}

func ToProductDTOs(products []*models.Product) []ProductDTO {
	var dtos []ProductDTO
	for _, product := range products {
		dtos = append(dtos, ToProductDTO(product))
	}
	return dtos
}
