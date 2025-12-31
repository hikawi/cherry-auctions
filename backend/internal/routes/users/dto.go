package users

import (
	"time"

	"luny.dev/cherryauctions/internal/models"
	"luny.dev/cherryauctions/pkg/ranges"
)

type SubscriptionDTO struct {
	ExpiredAt time.Time `json:"expired_at"`
	CreatedAt time.Time `json:"created_at"`
}

type UserDTO struct {
	ID              uint             `json:"id"`
	Name            string           `json:"name"`
	Email           *string          `json:"email"`
	Verified        bool             `json:"verified"`
	CreatedAt       time.Time        `json:"created_at"`
	AverageRating   float64          `json:"average_rating"`
	WaitingApproval bool             `json:"waiting_approval"`
	Roles           []string         `json:"roles"`
	Subscription    *SubscriptionDTO `json:"subscription"`
}

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
