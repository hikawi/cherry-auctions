package services

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"luny.dev/cherryauctions/internal/models"
	"luny.dev/cherryauctions/internal/repositories"
)

type OTPService struct {
	mailer   *MailerService
	userRepo *repositories.UserRepository
}

var (
	ErrOTPDidntUpdate = errors.New("couldn't update otp in user repo, wrong id?")
	ErrOTPWrongOTP    = errors.New("wrong otp")
	ErrOTPCantClear   = errors.New("couldn't clear otp")
)

func NewOTPService(
	mailer *MailerService,
	userRepo *repositories.UserRepository,
) *OTPService {
	return &OTPService{
		mailer:   mailer,
		userRepo: userRepo,
	}
}

func (s *OTPService) VerifyOTP(ctx context.Context, userID uint, otpCode string) error {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.OTPCode == nil || *user.OTPCode != otpCode || user.OTPExpiredAt == nil || user.OTPExpiredAt.Before(time.Now()) {
		return ErrOTPWrongOTP
	}

	rows, err := s.userRepo.ClearOTP(ctx, userID)
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrOTPCantClear
	}

	return nil
}

// SendOTP sends an email to the OTP.
func (s *OTPService) SendOTP(ctx context.Context, user *models.User) error {
	otp, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return err
	}

	otpCode := otp.Int64() + 100000
	rows, err := s.userRepo.UpdateOTP(ctx, user.ID, fmt.Sprintf("%d", otpCode))
	if err != nil {
		return err
	}

	if rows != 1 {
		return ErrOTPDidntUpdate
	}

	s.mailer.SendOTPEmail(user, fmt.Sprintf("%d", otpCode))
	return nil
}
