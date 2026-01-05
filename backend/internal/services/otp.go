package services

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"luny.dev/cherryauctions/internal/models"
	"luny.dev/cherryauctions/internal/repositories"
)

type OTPService struct {
	mailer   *MailerService
	userRepo *repositories.UserRepository
}

var ErrOTPDidntUpdate = errors.New("couldn't update otp in user repo, wrong id?")

func NewOTPService(
	mailer *MailerService,
	userRepo *repositories.UserRepository,
) *OTPService {
	return &OTPService{
		mailer:   mailer,
		userRepo: userRepo,
	}
}

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
