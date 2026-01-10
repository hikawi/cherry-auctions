package shared

import "time"

type ProfileDTO struct {
	ID            uint    `json:"id"`
	Name          *string `json:"name"`
	Email         *string `json:"email"`
	AvatarURL     *string `json:"avatar_url"`
	AverageRating float64 `json:"average_rating"`
}

type BidDTO struct {
	ID        uint       `json:"id"`
	Price     int64      `json:"price"`
	Automated bool       `json:"automated"`
	Bidder    ProfileDTO `json:"bidder"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type ProductDTO struct {
	ID                  uint       `json:"id"`
	Name                string     `json:"name"`
	StartingBid         int64      `json:"starting_bid"`
	StepBidValue        int64      `json:"step_bid_value"`
	BINPrice            *int64     `json:"bin_price"`
	Description         string     `json:"description"`
	ThumbnailURL        string     `json:"thumbnail_url"`
	AllowsUnratedBuyers bool       `json:"allows_unrated_buyers"`
	AutoExtendsTime     bool       `json:"auto_extends_time"`
	CreatedAt           time.Time  `json:"created_at"`
	ExpiredAt           time.Time  `json:"expired_at"`
	Seller              ProfileDTO `json:"seller"`
	// CurrentHighestBid   *BidDTO                `json:"current_highest_bid"`
	// Categories          []CategoryDTO          `json:"categories"`
	// DescriptionChanges  []DescriptionChangeDTO `json:"description_changes"`
	// DeniedBidders []ProfileDTO `json:"denied_bidders"`
	BidsCount    int        `json:"bids_count"`
	IsFavorite   bool       `json:"is_favorite"`
	ProductState string     `json:"product_state"`
	FinalizedAt  *time.Time `json:"finalized_at"`
}

type ChatSessionDTO struct {
	ID      uint       `json:"id"`
	Seller  ProfileDTO `json:"seller"`
	Buyer   ProfileDTO `json:"buyer"`
	Product ProductDTO `json:"product"`
}

type ChatMessageDTO struct {
	ID            uint       `json:"id"`
	Sender        ProfileDTO `json:"sender"`
	Content       string     `json:"content"`
	ImageURL      *string    `json:"image_url"`
	ChatSessionID uint       `json:"chat_session_id"`
}
