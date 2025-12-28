package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"luny.dev/cherryauctions/internal/services"
)

func TestPasswordHash(t *testing.T) {
	randomService := services.RandomService{}
	passwordService := services.PasswordService{RandomService: &randomService}

	t.Run("CantVerify", func(t *testing.T) {
		testHash := "lol"
		_, err := passwordService.VerifyPassword(testHash, "lol")
		assert.NotNil(t, err)
	})

	t.Run("Matches", func(t *testing.T) {
		testPassword := "test@123"
		hash, err := passwordService.HashPassword(testPassword)
		assert.Nil(t, err)

		ok, err := passwordService.VerifyPassword(hash, "test@123")
		assert.True(t, ok)
		assert.Nil(t, err)
	})

	t.Run("DoesntMatch", func(t *testing.T) {
		testPassword := "test@123"
		hash, err := passwordService.HashPassword(testPassword)
		assert.Nil(t, err)

		ok, err := passwordService.VerifyPassword(hash, "test@12")
		assert.False(t, ok)
		assert.Nil(t, err)
	})

	t.Run("HashDifferentForSameInput", func(t *testing.T) {
		testPassword := "test@123"
		hash1, err := passwordService.HashPassword(testPassword)
		assert.Nil(t, err)

		hash2, err := passwordService.HashPassword(testPassword)
		assert.Nil(t, err)

		assert.NotEqualValues(t, hash1, hash2)
	})
}
