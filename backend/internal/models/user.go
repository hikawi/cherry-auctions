// Package models provide a set of GORM-based modelling for the database.
package models

import (
	"time"
)

// User is the main struct that binds everything together.
type User struct {
	ID           uint    `gorm:"column:id;primaryKey;autoIncrement"`
	Name         *string `gorm:"column:name;size:200"`
	Email        *string `gorm:"column:email;size:200;unique;uniqueIndex"`
	Password     *string `gorm:"column:password"`
	Address      *string
	AvatarURL    *string
	OauthType    string     `gorm:"column:oauth_type;not null;default:none;check:oauth_type in ('google','none')"`
	Verified     bool       `gorm:"column:verified;not null;default:false"`
	OTPCode      *string    `gorm:"column:otp_code;size:10"`
	OTPExpiredAt *time.Time `gorm:"column:otp_expired_at"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime;not null"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime;not null"`

	AverageRating   float64 `gorm:"not null;default:0"`
	WaitingApproval bool    `gorm:"not null;default:false"`

	RefreshTokens    []RefreshToken       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Roles            []Role               `gorm:"many2many:user_roles"`
	Subscriptions    []SellerSubscription `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	FavoriteProducts []Product            `gorm:"many2many:favorite_products"`
	Ratings          []Rating             `gorm:"foreignKey:ReviewerID"`
	RatedRatings     []Rating             `gorm:"foreignKey:RevieweeID"`
}
