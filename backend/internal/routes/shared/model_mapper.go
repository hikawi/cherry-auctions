package shared

import (
	"luny.dev/cherryauctions/internal/models"
)

func ToProfileDTO(m *models.User) ProfileDTO {
	if m == nil {
		return ProfileDTO{}
	}

	return ProfileDTO{
		ID:            m.ID,
		Name:          m.Name,
		Email:         m.Email,
		AvatarURL:     m.AvatarURL,
		AverageRating: m.AverageRating,
	}
}

func ToBidDTO(m *models.Bid) BidDTO {
	if m == nil {
		return BidDTO{}
	}

	return BidDTO{
		ID:        m.ID,
		Price:     m.Price,
		Automated: m.Automated,
		Bidder:    ToProfileDTO(&m.User),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func ToChatSessionDTO(m *models.ChatSession) ChatSessionDTO {
	if m == nil {
		return ChatSessionDTO{}
	}

	return ChatSessionDTO{
		ID:      m.ID,
		Buyer:   ToProfileDTO(&m.Buyer),
		Seller:  ToProfileDTO(&m.Seller),
		Product: ToProductDTO(&m.Product),
	}
}

func ToChatMessageDTO(m *models.ChatMessage) ChatMessageDTO {
	if m == nil {
		return ChatMessageDTO{}
	}

	return ChatMessageDTO{
		ID:            m.ID,
		Sender:        ToProfileDTO(&m.Sender),
		Content:       m.Content,
		ImageURL:      m.ImageURL,
		ChatSessionID: m.ChatSessionID,
	}
}

func ToProductDTO(m *models.Product) ProductDTO {
	if m == nil {
		return ProductDTO{}
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
		Seller:              ToProfileDTO(&m.Seller),
		// CurrentHighestBid:   highestBid,
		BidsCount: m.BidsCount,
		// Categories:          ranges.Each(m.Categories, ToCategoryDTO),
		ProductState: string(m.ProductState),
		// DeniedBidders: ranges.Each(m.DeniedBidders, func(bidder models.DeniedBidder) ProfileDTO {
		// 	return ToProfileDTO(bidder.User)
		// }),
		// DescriptionChanges: ranges.Each(m.DescriptionChanges, func(m models.DescriptionChange) DescriptionChangeDTO {
		// 	return DescriptionChangeDTO{
		// 		ID:        m.ID,
		// 		Changes:   m.Changes,
		// 		CreatedAt: m.CreatedAt,
		// 	}
		// }),
		IsFavorite:  m.IsFavorite,
		FinalizedAt: m.FinalizedAt,
	}
}
