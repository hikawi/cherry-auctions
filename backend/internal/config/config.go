// Package config provides a centralized way to handle environment variables
package config

import "luny.dev/cherryauctions/pkg/env"

// -------------------------------------
// Thanks Gemini
// -------------------------------------

type Config struct {
	DatabaseURL     string
	Domain          string
	CookieSecure    bool
	RecaptchaSecret string

	CORS struct {
		Origins string
		Methods string
		Headers string
	}

	JWT struct {
		Secret   string
		Audience string
		Expiry   int
	}

	AWS struct {
		AccessKeyID     string
		SecretAccessKey string
		SessionToken    string
		S3Base          string
		S3UsePathStyle  bool
		BucketName      string
	}

	SMTP struct {
		Host     string
		Port     int
		User     string
		Password string
	}
}

func Load() *Config {
	cfg := &Config{}

	// General
	cfg.DatabaseURL = env.Fatalenv("DATABASE_URL")
	cfg.Domain = env.Fatalenv("DOMAIN")
	cfg.CookieSecure = env.FatalenvBool("COOKIE_SECURE")
	cfg.RecaptchaSecret = env.Fatalenv("RECAPTCHA_SECRET")

	// CORS
	cfg.CORS.Origins = env.Fatalenv("CORS_ORIGINS")
	cfg.CORS.Methods = env.Fatalenv("CORS_METHODS")
	cfg.CORS.Headers = env.Fatalenv("CORS_HEADERS")

	// JWT
	cfg.JWT.Secret = env.Fatalenv("JWT_SECRET_KEY")
	cfg.JWT.Audience = env.Fatalenv("JWT_AUDIENCE")
	cfg.JWT.Expiry = int(env.FatalenvInt("JWT_EXPIRY"))

	// AWS
	cfg.AWS.AccessKeyID = env.Fatalenv("AWS_ACCESS_KEY_ID")
	cfg.AWS.SecretAccessKey = env.Fatalenv("AWS_SECRET_ACCESS_KEY")
	cfg.AWS.SessionToken = env.Fatalenv("AWS_SESSION_TOKEN")
	cfg.AWS.S3Base = env.Fatalenv("AWS_S3_BASE")
	cfg.AWS.S3UsePathStyle = env.FatalenvBool("AWS_S3_USE_PATH_STYLE")
	cfg.AWS.BucketName = env.Fatalenv("AWS_BUCKET_NAME")

	// SMTP
	cfg.SMTP.Host = env.Fatalenv("SMTP_HOST")
	cfg.SMTP.Port = int(env.FatalenvInt("SMTP_PORT"))
	cfg.SMTP.User = env.Fatalenv("SMTP_USER")
	cfg.SMTP.Password = env.Fatalenv("SMTP_PASSWORD")

	return cfg
}
