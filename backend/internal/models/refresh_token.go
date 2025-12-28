package models

import "time"

type RefreshToken struct {
	ID           uint      `gorm:"column:id;primaryKey;autoIncrement"`
	UserID       uint      `gorm:"column:user_id;index"`
	User         User      `gorm:""`
	RefreshToken string    `gorm:"column:refresh_token;not null;uniqueIndex"`
	IsRevoked    bool      `gorm:"column:is_revoked;not null;default:false"`
	ExpiredAt    time.Time `gorm:"column:expired_at;not null"`
	CreatedAt    time.Time `gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;not null;autoUpdateTime"`
}
