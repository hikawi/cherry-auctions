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
	err := r.db.
		WithContext(ctx).
		Model(&models.Rating{}).
		Preload("Reviewer").
		Preload("Reviewee").
		Preload("Product").
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
	err := r.db.
		WithContext(ctx).
		Model(&models.Rating{}).
		Where("reviewer_id = ?", userID).
		Count(&count).
		Error
	return count, err
}

// GetMyReviewedRatings retrieves a user's list of ratings (they are the reviewee)
func (r *RatingRepostory) GetMyReviewedRatings(ctx context.Context, userID uint, limit int, offset int) ([]models.Rating, error) {
	var ratings []models.Rating
	err := r.db.
		WithContext(ctx).
		Model(&models.Rating{}).
		Preload("Reviewer").
		Preload("Reviewee").
		Preload("Product").
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
	err := r.db.
		WithContext(ctx).
		Model(&models.Rating{}).
		Where("reviewee_id = ?", userID).
		Count(&count).
		Error
	return count, err
}

func (r *RatingRepostory) GetRatingByID(ctx context.Context, id uint) (models.Rating, error) {
	var rating models.Rating
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&rating).
		Error
	return rating, err
}

// CreateRating creates a new rating
func (r *RatingRepostory) CreateRating(ctx context.Context, rating *models.Rating) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&models.Rating{}).Create(rating).Error
		if err != nil {
			return err
		}

		return tx.
			Model(&models.User{}).
			Where(rating.RevieweeID).
			Update("average_rating", tx.Model(&models.Rating{}).Select("avg(rating)").Where("reviewee_id = ?", rating.RevieweeID)).
			Error
	})
}

// UpdateRating updates an existing rating.
func (r *RatingRepostory) UpdateRating(ctx context.Context, ratingID uint, newRating uint, newFeedback string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var rating models.Rating
		err := tx.Model(&models.Rating{}).
			Select("id = ?", ratingID).
			First(&rating).
			Error
		if err != nil {
			return err
		}

		err = tx.Model(&rating).
			Select("rating", "feedback").
			UpdateColumns(map[string]any{
				"rating":   newRating,
				"feedback": newFeedback,
			}).
			Error
		if err != nil {
			return err
		}

		return tx.
			Model(&models.User{}).
			Where(rating.RevieweeID).
			Update("average_rating", tx.Model(&models.Rating{}).Select("avg(rating)").Where("reviewee_id = ?", rating.RevieweeID)).
			Error
	})
}

// DeleteRatingByID deletes a rating by ID.
func (r *RatingRepostory) DeleteRatingByID(ctx context.Context, ratingID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var rating models.Rating
		err := tx.Model(&models.Rating{}).
			Where("id = ?", ratingID).
			First(&rating).
			Error
		if err != nil {
			return err
		}

		err = tx.Model(&rating).
			Where("id = ?", ratingID).
			Delete(&rating).
			Error
		if err != nil {
			return err
		}

		return tx.Exec("UPDATE users SET average_rating = (SELECT AVG(rating) FROM ratings WHERE reviewee_id = users.id) WHERE id IN (?, ?)", rating.RevieweeID, rating.ReviewerID).Error
	})
}
