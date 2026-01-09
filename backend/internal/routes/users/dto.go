package users

import (
	"mime/multipart"
	"time"
)

type QuestionDTO struct {
	ID        uint       `json:"id"`
	Content   string     `json:"content"`
	Answer    *string    `json:"answer"`
	User      ProfileDTO `json:"user"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type SubscriptionDTO struct {
	ExpiredAt time.Time `json:"expired_at"`
	CreatedAt time.Time `json:"created_at"`
}

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
	Bids                []BidDTO               `json:"bids"`

	BidsCount  int  `json:"bids_count"`
	IsFavorite bool `json:"is_favorite"`
}

type GetMyProductsQuery struct {
	Page    int `form:"page" binding:"number,gt=0,omitempty" json:"page"`
	PerPage int `form:"per_page" binding:"number,gt=0,omitempty" json:"per_page"`
}

type GetProductsResponse struct {
	Data       []ProductDTO `json:"data"`
	Total      int64        `json:"total"`
	TotalPages int          `json:"total_pages"`
	Page       int          `json:"page"`
	PerPage    int          `json:"per_page"`
}

type UserDTO struct {
	ID              uint             `json:"id"`
	Name            *string          `json:"name"`
	Email           *string          `json:"email"`
	Address         *string          `json:"address"`
	AvatarURL       *string          `json:"avatar_url"`
	Verified        bool             `json:"verified"`
	CreatedAt       time.Time        `json:"created_at"`
	AverageRating   float64          `json:"average_rating"`
	WaitingApproval bool             `json:"waiting_approval"`
	Roles           []string         `json:"roles"`
	Subscription    *SubscriptionDTO `json:"subscription"`
}

type GetUsersQuery struct {
	Query   string `form:"query" json:"query"`
	Page    int    `form:"page" binding:"number,gt=0,omitempty" json:"page"`
	PerPage int    `form:"per_page" binding:"number,gt=0,omitempty" json:"per_page"`
}

type GetUsersResponse struct {
	Data       []UserDTO `json:"data"`
	Total      int64     `json:"total"`
	TotalPages int       `json:"total_pages"`
	Page       int       `json:"page"`
	PerPage    int       `json:"per_page"`
}

type PostApproveRequest struct {
	ID int `json:"id" binding:"number,gt=0"`
}

type PostAvatarRequest struct {
	Avatar *multipart.FileHeader `form:"avatar" binding:"required"`
}

type PostAvatarResponse struct {
	AvatarURL string `json:"avatar_url"`
}

type PostProfileRequest struct {
	Name    *string `json:"name" form:"name" binding:"min=2,max=200,omitempty"`
	Address *string `json:"address" form:"address" binding:"min=2,omitempty"`
}

type PutPasswordRequest struct {
	NewPassword     string `json:"new_password" form:"new_password" binding:"min=2,required"`
	CurrentPassword string `json:"current_password" form:"current_password" binding:"min=2,required"`
}
