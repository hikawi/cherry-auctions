package models

import (
	"time"

	"gorm.io/gorm"
)

type BidIntent struct {
	ProductID uint           `gorm:"primaryKey"`
	UserID    uint           `gorm:"primaryKey"`
	BidAmount int64          `gorm:"not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
