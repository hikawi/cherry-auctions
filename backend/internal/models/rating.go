package models

import "gorm.io/gorm"

type Rating struct {
	gorm.Model
	Rating     float64 `gorm:"not null"`
	ReviewerID uint
	Reviewer   User
	RevieweeID uint
	Reviewee   User
}

func (r *Rating) AfterSave(tx *gorm.DB) error {
	var avgRating float64

	err := tx.Model(&Rating{}).
		Where("reviewee_id = ?", r.RevieweeID).
		Select("COALESCE(AVG(rating), 0)").
		Row().
		Scan(&avgRating)
	if err != nil {
		return err
	}

	return tx.Model(&User{}).
		Where("id = ?", r.RevieweeID).
		Update("average_rating", avgRating).Error
}

func (r *Rating) AfterDelete(tx *gorm.DB) error {
	return r.AfterSave(tx)
}
