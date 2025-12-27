package services

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"luny.dev/cherryauctions/utils"
)

type JWTSubject struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// SignJWT signs a JWT based on the environment variables and returns a signed string.
func SignJWT(id uint, email string, role string) (string, error) {
	expiryTime, err := strconv.ParseUint(utils.Fatalenv("JWT_EXPIRY"), 10, 64)

	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTSubject{
		id,
		email,
		role,
		jwt.RegisteredClaims{
			Issuer:    utils.Fatalenv("DOMAIN"),
			Audience:  jwt.ClaimStrings{utils.Fatalenv("JWT_AUDIENCE")},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   email,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryTime) * time.Second)),
		},
	})
	str, err := token.SignedString([]byte(utils.Fatalenv("JWT_SECRET_KEY")))
	return str, err
}

// VerifyJWT verifies if a JWT is valid under some conditions.
func VerifyJWT(signedString string) (*JWTSubject, error) {
	var sub JWTSubject

	parser := jwt.NewParser(jwt.WithAudience(utils.Fatalenv("JWT_AUDIENCE")),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(utils.Fatalenv("DOMAIN")))
	token, err := parser.ParseWithClaims(signedString, &sub, func(t *jwt.Token) (any, error) {
		return []byte(utils.Fatalenv("JWT_SECRET_KEY")), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return &sub, err
}
