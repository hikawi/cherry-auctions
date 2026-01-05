package users

import (
	"luny.dev/cherryauctions/internal/models"
	"luny.dev/cherryauctions/pkg/ranges"
)

func ToSubscriptionDTO(m models.SellerSubscription) SubscriptionDTO {
	return SubscriptionDTO{
		ExpiredAt: m.ExpiredAt,
		CreatedAt: m.CreatedAt,
	}
}

func ToUserDTO(m *models.User) UserDTO {
	var subscription *SubscriptionDTO
	if len(m.Subscriptions) > 0 {
		dto := ToSubscriptionDTO(m.Subscriptions[0])
		subscription = &dto
	}

	return UserDTO{
		ID:              m.ID,
		Name:            m.Name,
		Email:           m.Email,
		Address:         m.Address,
		AvatarURL:       m.AvatarURL,
		Verified:        m.Verified,
		CreatedAt:       m.CreatedAt,
		AverageRating:   m.AverageRating,
		WaitingApproval: m.WaitingApproval,
		Roles: ranges.Each(m.Roles, func(r models.Role) string {
			return r.ID
		}),
		Subscription: subscription,
	}
}

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

func ToProductImageDTO(m models.ProductImage) ProductImageDTO {
	return ProductImageDTO{
		URL:     m.URL,
		AltText: m.AltText,
	}
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
		Bids:                ranges.Each(m.Bids, ToBidDTO),
		Categories:          ranges.Each(m.Categories, ToCategoryDTO),
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
