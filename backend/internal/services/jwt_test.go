package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"luny.dev/cherryauctions/internal/services"
)

func TestSignJWT(t *testing.T) {
	jwtService := services.JWTService{
		JWTExpiry:    30,
		JWTAudience:  "test",
		JWTDomain:    "https://example.com",
		JWTSecretKey: "test",
	}
	str, err := jwtService.SignJWT(2, "test@example.com", "lol hi")

	assert.Nil(t, err)
	assert.NotNil(t, str)

	// Revalidate
	t.Run("SuccessValidation", func(t *testing.T) {
		claims, err := jwtService.VerifyJWT(str)
		assert.Nil(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, claims.ID, uint(2))
		assert.Equal(t, claims.Email, "test@example.com")
		assert.Equal(t, claims.Roles, "lol hi")
	})

	jwtService2 := services.JWTService{
		JWTExpiry:    30,
		JWTAudience:  "test",
		JWTDomain:    "https://example.com",
		JWTSecretKey: "test2",
	}

	// Revalidate but with different secret key or wrong key.
	t.Run("FailedValidation", func(t *testing.T) {
		_, err := jwtService2.VerifyJWT(str)
		assert.NotNil(t, err)
	})
}
