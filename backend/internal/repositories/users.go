package repositories

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
	"luny.dev/cherryauctions/internal/models"
)

type UserRepository struct {
	DB             *gorm.DB
	RoleRepository *RoleRepository
}

// GetUserByID retrieves a single user using an ID.
func (repo *UserRepository) GetUserByID(ctx context.Context, id uint) (models.User, error) {
	return gorm.G[models.User](repo.DB).
		Where("id = ?", id).
		Preload("Roles", nil).
		Preload("Subscriptions", func(db gorm.PreloadBuilder) error {
			db.Where("expired_at > ?", time.Now()).Limit(1)
			return nil
		}).
		First(ctx)
}

// GetUserByEmail returns a user with the email, if it found.
// An error is returned if the user can not be found.
// Email is insensitive.
func (repo *UserRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	return gorm.G[models.User](repo.DB).
		Preload("Roles", nil).
		Preload("Subscriptions", func(db gorm.PreloadBuilder) error {
			db.Where("expired_at > ?", time.Now()).Limit(1)
			return nil
		}).
		Where("email ILIKE ?", strings.ToLower(email)).
		First(ctx)
}

func (repo *UserRepository) GetUsers(ctx context.Context, limit int, offset int) ([]models.User, error) {
	return gorm.G[models.User](repo.DB).
		Preload("Roles", nil).
		Preload("Subscriptions", func(db gorm.PreloadBuilder) error {
			db.Where("expired_at > ?", time.Now()).Limit(1)
			return nil
		}).
		Order("id").
		Limit(limit).
		Offset(offset).
		Find(ctx)
}

// RegisterNewUser registers a new user with a default role.
func (repo *UserRepository) RegisterNewUser(ctx context.Context, name string, email string, password string) (models.User, error) {
	defaultRole, err := repo.RoleRepository.GetRoleByID(ctx, "user")
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		Name:      &name,
		Email:     &email,
		Password:  &password,
		OauthType: "none",
		Roles:     []models.Role{defaultRole},
	}
	err = gorm.G[models.User](repo.DB).Create(ctx, &user)
	return user, err
}

func (repo *UserRepository) CountUsers(ctx context.Context) (int64, error) {
	return gorm.G[models.User](repo.DB).Count(ctx, "id")
}

// SaveUser creates a new user with the model passed in.
// Returns an error if it can't be saved.
func (repo *UserRepository) SaveUser(ctx context.Context, user *models.User) error {
	return gorm.G[models.User](repo.DB).Create(ctx, user)
}

// RequestUserApproval marks a user as requesting approval.
func (repo *UserRepository) RequestUserApproval(ctx context.Context, id uint) (int, error) {
	return gorm.G[models.User](repo.DB).Where("id = ?", id).Update(ctx, "waiting_approval", true)
}

func (repo *UserRepository) ApproveUser(ctx context.Context, id uint) error {
	return repo.DB.Transaction(func(tx *gorm.DB) error {
		rows, err := gorm.G[models.User](tx).Where("id = ?", id).Update(ctx, "waiting_approval", false)
		if err != nil {
			return err
		}
		if rows != 1 {
			return errors.New("couldn't mark as no longer waiting for approval")
		}

		// Add a subscription
		subscription := models.SellerSubscription{
			UserID:    id,
			ExpiredAt: time.Now().Add(time.Minute * 60 * 24 * 7), // 7 days, should be made configured later
		}
		err = gorm.G[models.SellerSubscription](tx).Create(ctx, &subscription)
		if err != nil {
			return err
		}

		return nil
	})
}

func (repo *UserRepository) UpdateAvatarURL(ctx context.Context, id uint, url string) (int, error) {
	return gorm.G[models.User](repo.DB).Where("id = ?", id).Update(ctx, "avatar_url", url)
}

func (repo *UserRepository) UpdateProfile(ctx context.Context, id uint, name *string, address *string) (int, error) {
	return gorm.G[models.User](repo.DB).Where("id = ?", id).Updates(ctx, models.User{Name: name, Address: address})
}

func (repo *UserRepository) UpdatePassword(ctx context.Context, id uint, password string) (int, error) {
	return gorm.G[models.User](repo.DB).Where("id = ?", id).Update(ctx, "password", password)
}

// UpdateOTP updates the user's OTP to a new one.
func (repo *UserRepository) UpdateOTP(ctx context.Context, id uint, otp string) (int, error) {
	expiredAt := time.Now().Add(15 * time.Minute)
	return gorm.G[models.User](repo.DB).
		Where("id = ?", id).
		Updates(ctx, models.User{OTPCode: &otp, OTPExpiredAt: &expiredAt})
}

// ClearOTP clears the user's OTP to an empty state to mark used.
func (repo *UserRepository) ClearOTP(ctx context.Context, id uint) (int, error) {
	db := repo.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Select("otp_code", "otp_expired_at").
		Updates(map[string]any{"otp_code": nil, "otp_expired_at": nil})
	return int(db.RowsAffected), db.Error
}

func (repo *UserRepository) UpdateUserVerified(ctx context.Context, id uint, verified bool) (int, error) {
	return gorm.G[models.User](repo.DB).Where("id = ?", id).Update(ctx, "verified", verified)
}
