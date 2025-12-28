package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	JWTExpiry    int
	JWTDomain    string
	JWTAudience  string
	JWTSecretKey string
}

type JWTSubject struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Roles  string `json:"roles"`
	jwt.RegisteredClaims
}

// SignJWT signs a JWT based on the environment variables and returns a signed string.
func (s *JWTService) SignJWT(id uint, email string, roles string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTSubject{
		id,
		email,
		roles,
		jwt.RegisteredClaims{
			Issuer:    s.JWTDomain,
			Audience:  jwt.ClaimStrings{s.JWTAudience},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   email,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.JWTExpiry) * time.Second)),
		},
	})
	str, err := token.SignedString([]byte(s.JWTSecretKey))
	return str, err
}

// VerifyJWT verifies if a JWT is valid under some conditions.
func (s *JWTService) VerifyJWT(signedString string) (*JWTSubject, error) {
	var sub JWTSubject

	parser := jwt.NewParser(jwt.WithAudience(s.JWTAudience),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(s.JWTDomain))
	token, err := parser.ParseWithClaims(signedString, &sub, func(t *jwt.Token) (any, error) {
		return []byte(s.JWTSecretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return &sub, err
}
