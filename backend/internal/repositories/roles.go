package repositories

import (
	"context"

	"gorm.io/gorm"
	"luny.dev/cherryauctions/internal/models"
)

type RoleRepository struct {
	DB *gorm.DB
}

func (r *RoleRepository) GetRoles(ctx context.Context) ([]models.Role, error) {
	return gorm.G[models.Role](r.DB).Find(ctx)
}

func (r *RoleRepository) GetRoleByID(ctx context.Context, id string) (models.Role, error) {
	return gorm.G[models.Role](r.DB).Where("id = ?", id).First(ctx)
}

func (r *RoleRepository) SaveRole(ctx context.Context, role *models.Role) error {
	return gorm.G[models.Role](r.DB).Create(ctx, role)
}
