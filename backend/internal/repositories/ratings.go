package repositories

import (
	"context"

	"gorm.io/gorm"
	"luny.dev/cherryauctions/internal/models"
)

type RatingRepostory struct {
	db *gorm.DB
}

func NewRatingRepository(
	db *gorm.DB,
) *RatingRepostory {
	return &RatingRepostory{
		db: db,
	}
}

// GetMyRatings retrieves a user's list of ratings (they are the reviewer)
func (r *RatingRepostory) GetMyRatings(ctx context.Context, userID uint, limit int, offset int) ([]models.Rating, error) {
	var ratings []models.Rating
	err := r.db.Model(&models.Rating{}).
		Preload("Reviewer").
		Preload("Reviewee").
		Where("reviewer_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&ratings).
		Error
	return ratings, err
}

func (r *RatingRepostory) CountMyRatings(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Rating{}).
		Where("reviewer_id = ?", userID).
		Count(&count).
		Error
	return count, err
}

// GetMyReviewedRatings retrieves a user's list of ratings (they are the reviewee)
func (r *RatingRepostory) GetMyReviewedRatings(ctx context.Context, userID uint, limit int, offset int) ([]models.Rating, error) {
	var ratings []models.Rating
	err := r.db.Model(&models.Rating{}).
		Preload("Reviewer").
		Preload("Reviewee").
		Where("reviewee_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&ratings).
		Error
	return ratings, err
}

func (r *RatingRepostory) CountMyReviewedRatings(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Rating{}).
		Where("reviewee_id = ?", userID).
		Count(&count).
		Error
	return count, err
}
