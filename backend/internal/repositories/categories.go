package repositories

import (
	"context"

	"gorm.io/gorm"
	"luny.dev/cherryauctions/internal/models"
)

type CategoryRepository struct {
	DB *gorm.DB
}

// GetActiveCategories retrieves a list of categories that are not deleted yet.
// This returns a flat-list, that does not populate anything.
func (repo *CategoryRepository) GetActiveCategories(ctx context.Context) ([]models.Category, error) {
	return gorm.G[models.Category](repo.DB).Find(ctx)
}

// SaveCategory saves the category in the database, returns an error if it couldn't.
func (repo *CategoryRepository) SaveCategory(ctx context.Context, category *models.Category) error {
	return gorm.G[models.Category](repo.DB).Create(ctx, category)
}
