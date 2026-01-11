package repositories

import (
	"context"

	"gorm.io/gorm"
	"luny.dev/cherryauctions/internal/models"
)

type BidIntentRepository struct {
	db *gorm.DB
}

func NewBidIntentRepository(
	db *gorm.DB,
) *BidIntentRepository {
	return &BidIntentRepository{
		db: db,
	}
}

func (r *BidIntentRepository) CreateBidIntent(ctx context.Context, bidIntent *models.BidIntent) error {
	// upsert
	return r.db.WithContext(ctx).
		Model(&models.BidIntent{}).
		Save(&bidIntent).
		Error
}
