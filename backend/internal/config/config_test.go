package config_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"luny.dev/cherryauctions/internal/config"
)

func TestConfigLoad(t *testing.T) {
	// 1. Define the "Expected" values
	expected := struct {
		dbURL    string
		jwtSec   string
		smtpPort int
		s3Path   bool
	}{
		dbURL:    "postgres://user:pass@localhost:5432/db",
		jwtSec:   "super-secret-key",
		smtpPort: 587,
		s3Path:   true,
	}

	// 2. Set the environment variables
	// Note: You must set EVERY variable your Fatalenv calls expect
	t.Setenv("DATABASE_URL", expected.dbURL)
	t.Setenv("DOMAIN", "localhost")
	t.Setenv("COOKIE_SECURE", "true")
	t.Setenv("RECAPTCHA_SECRET", "captcha-key")

	t.Setenv("CORS_ORIGINS", "*")
	t.Setenv("CORS_METHODS", "GET,POST")
	t.Setenv("CORS_HEADERS", "Content-Type")

	t.Setenv("JWT_SECRET_KEY", expected.jwtSec)
	t.Setenv("JWT_AUDIENCE", "cherry-auctions")
	t.Setenv("JWT_EXPIRY", "3600")

	t.Setenv("AWS_ACCESS_KEY_ID", "AKIA...")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
	t.Setenv("AWS_SESSION_TOKEN", "token")
	t.Setenv("AWS_S3_BASE", "s3.amazonaws.com")
	t.Setenv("AWS_S3_USE_PATH_STYLE", strconv.FormatBool(expected.s3Path))
	t.Setenv("AWS_BUCKET_NAME", "my-bucket")

	t.Setenv("SMTP_HOST", "smtp.gmail.com")
	t.Setenv("SMTP_PORT", strconv.Itoa(expected.smtpPort))
	t.Setenv("SMTP_USER", "user@gmail.com")
	t.Setenv("SMTP_PASSWORD", "password")

	// 3. Run the function
	cfg := config.Load()

	// 4. Assertions
	assert.EqualValues(t, expected.dbURL, cfg.DatabaseURL)
	assert.EqualValues(t, expected.jwtSec, cfg.JWT.Secret)
	assert.EqualValues(t, expected.smtpPort, cfg.SMTP.Port)
	assert.EqualValues(t, expected.s3Path, cfg.AWS.S3UsePathStyle)
}
