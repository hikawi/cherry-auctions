package repositories

import (
	"context"
	"strings"

	"gorm.io/gorm"
	"luny.dev/cherryauctions/models"
)

type UserRepository struct {
	DB *gorm.DB
}

// GetUserByEmail returns a user with the email, if it found.
// An error is returned if the user can not be found.
// Email is insensitive.
func (repo *UserRepository) GetUserByEmail(email string) (models.User, error) {
	ctx := context.Background()
	return gorm.G[models.User](repo.DB).Where("email ILIKE ?", strings.ToLower(email)).First(ctx)
}

// SaveUser creates a new user with the model passed in.
// Returns an error if it can't be saved.
func (repo *UserRepository) SaveUser(user *models.User) error {
	ctx := context.Background()
	return gorm.G[models.User](repo.DB).Create(ctx, user)
}
