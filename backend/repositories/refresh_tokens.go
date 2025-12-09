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
		ExpiredAt:    time.Now().Add(time.Hour * 24 * 30),
		IsRevoked:    false,
	}
	err := gorm.G[models.RefreshToken](repo.DB).Create(ctx, &refreshToken)
	return refreshToken, err
}

func (repo *RefreshTokenRepository) GetUserToken(id uint, token string) (models.RefreshToken, error) {
	ctx := context.Background()
	refreshToken, err := gorm.G[models.RefreshToken](repo.DB).Where("user_id = ? AND token = ?", id, token).First(ctx)
	return refreshToken, err
}

func (repo *RefreshTokenRepository) InvalidateToken(token string) (int, error) {
	ctx := context.Background()
	return gorm.G[models.RefreshToken](repo.DB).Where("token = ?", token).Update(ctx, "is_revoked", true)
}
