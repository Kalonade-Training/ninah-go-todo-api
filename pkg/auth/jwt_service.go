package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	Generate(userID string) (string, error)
}

type jwtService struct{ secret string }

func NewJWTService(secret string) TokenService {
	return &jwtService{secret: secret}
}

func (s *jwtService) Generate(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(8 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(s.secret))
}
