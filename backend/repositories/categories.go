package repositories

import (
	"context"

	"gorm.io/gorm"
	"luny.dev/cherryauctions/models"
)

type CategoryRepository struct {
	DB      *gorm.DB
	Context context.Context
}

func (repo *CategoryRepository) GetActiveCategories() ([]models.Category, error) {
	return gorm.G[models.Category](repo.DB).Find(repo.Context)
}
