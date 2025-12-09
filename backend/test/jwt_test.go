package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"luny.dev/cherryauctions/services"
)

func TestSignJWT(t *testing.T) {
	t.Setenv("JWT_EXPIRY", "30")
	t.Setenv("JWT_AUDIENCE", "test")
	t.Setenv("DOMAIN", "https://example.com")
	t.Setenv("JWT_SECRET_KEY", "test")

	str, err := services.SignJWT(2, "test@example.com")

	assert.Nil(t, err)
	assert.NotNil(t, str)

	// Revalidate
	t.Run("SuccessValidation", func(t *testing.T) {
		claims, err := services.VerifyJWT(str)
		assert.Nil(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, claims.ID, uint(2))
		assert.Equal(t, claims.Email, "test@example.com")
	})

	// Revalidate but with different secret key or wrong key.
	t.Run("FailedValidation", func(t *testing.T) {
		t.Setenv("JWT_SECRET_KEY", "test2")

		_, err := services.VerifyJWT(str)
		assert.NotNil(t, err)
	})
}
