package repositories

import (
	"context"
	"strings"

	"gorm.io/gorm"
	"luny.dev/cherryauctions/internal/models"
)

type UserRepository struct {
	DB             *gorm.DB
	RoleRepository *RoleRepository
}

// GetUserByID retrieves a single user using an ID.
func (repo *UserRepository) GetUserByID(ctx context.Context, id uint) (models.User, error) {
	return gorm.G[models.User](repo.DB).Where("id = ?", id).Preload("Roles", nil).First(ctx)
}

// GetUserByEmail returns a user with the email, if it found.
// An error is returned if the user can not be found.
// Email is insensitive.
func (repo *UserRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	return gorm.G[models.User](repo.DB).Preload("Roles", nil).Where("email ILIKE ?", strings.ToLower(email)).First(ctx)
}

// RegisterNewUser registers a new user with a default role.
func (repo *UserRepository) RegisterNewUser(ctx context.Context, name string, email string, password string) (models.User, error) {
	defaultRole, err := repo.RoleRepository.GetRoleByID(ctx, "user")
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		Name:      name,
		Email:     &email,
		Password:  &password,
		OauthType: "none",
		Roles:     []models.Role{defaultRole},
	}
	err = gorm.G[models.User](repo.DB).Create(ctx, &user)
	return user, err
}

// SaveUser creates a new user with the model passed in.
// Returns an error if it can't be saved.
func (repo *UserRepository) SaveUser(ctx context.Context, user *models.User) error {
	return gorm.G[models.User](repo.DB).Create(ctx, user)
}
