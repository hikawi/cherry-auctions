package repositories

import (
	"context"
	"time"

	"gorm.io/gorm"
	"luny.dev/cherryauctions/models"
)

type RefreshTokenRepository struct {
	DB *gorm.DB
}

// SaveUserToken saves the token mapped to the user_id provided.
// This function does not do the hashing, do it beforehand before passing into this function.
func (repo *RefreshTokenRepository) SaveUserToken(id uint, token string) (models.RefreshToken, error) {
	ctx := context.Background()

	refreshToken := models.RefreshToken{
		UserID:       id,
		RefreshToken: token,
		ExpiredAt:    time.Now().Add(time.Hour * 24 * 30 * 3),
		IsRevoked:    false,
	}
	err := gorm.G[models.RefreshToken](repo.DB).Create(ctx, &refreshToken)
	return refreshToken, err
}

func (repo *RefreshTokenRepository) GetRefreshToken(token string) (models.RefreshToken, error) {
	ctx := context.Background()
	refreshToken, err := gorm.G[models.RefreshToken](repo.DB).Preload("User", nil).Where("refresh_token = ?", token).First(ctx)
	return refreshToken, err
}

// InvalidateToken invalidates a token by marking it as revoked.
// This function does not hash the token before checking.
func (repo *RefreshTokenRepository) InvalidateToken(token string) (int, error) {
	ctx := context.Background()
	return gorm.G[models.RefreshToken](repo.DB).Where("refresh_token = ?", token).Update(ctx, "is_revoked", true)
}
