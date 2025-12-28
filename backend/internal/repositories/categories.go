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
	return gorm.G[models.Category](repo.DB).Order("id asc").Find(ctx)
}

func (repo *CategoryRepository) GetCategoryByID(ctx context.Context, id uint) (models.Category, error) {
	return gorm.G[models.Category](repo.DB).Where("id = ?", id).First(ctx)
}

// SaveCategory saves the category in the database, returns an error if it couldn't.
func (repo *CategoryRepository) SaveCategory(ctx context.Context, category *models.Category) error {
	return gorm.G[models.Category](repo.DB).Create(ctx, category)
}

// UpdateCategory updates a category partially.
func (repo *CategoryRepository) UpdateCategory(ctx context.Context, id uint, name *string, parentID *uint) (int, error) {
	newCat := models.Category{}
	if name != nil {
		newCat.Name = *name
	}
	newCat.ParentID = parentID
	return gorm.G[models.Category](repo.DB).Select("Name", "ParentID").Omit("id").Where("id = ?", id).Updates(ctx, newCat)
}

func (repo *CategoryRepository) SoftDeleteCategory(ctx context.Context, id uint) (int, error) {
	return gorm.G[models.Category](repo.DB).Where("id = ?", id).Delete(ctx)
}
