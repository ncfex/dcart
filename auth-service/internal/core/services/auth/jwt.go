package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ncfex/dcart/auth-service/internal/domain/errors"
)

type JWTService struct {
	issuer      string
	tokenSecret string
}

func NewJWTService(issuer, tokenSecret string) *JWTService {
	return &JWTService{
		issuer:      issuer,
		tokenSecret: tokenSecret,
	}
}

func (s *JWTService) MakeJWT(userID uuid.UUID, expiresIn time.Duration) (string, error) {
	currentTime := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    s.issuer,
		IssuedAt:  jwt.NewNumericDate(currentTime.UTC()),
		ExpiresAt: jwt.NewNumericDate(currentTime.Add(expiresIn)),
		Subject:   userID.String(),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.tokenSecret))
}

func (s *JWTService) ValidateJWT(tokenString string) (uuid.UUID, error) {
	// check issuer
	// add more checks
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(t *jwt.Token) (interface{}, error) { return []byte(s.tokenSecret), nil },
	)
	if err != nil {
		return uuid.UUID{}, err
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return uuid.UUID{}, errors.ErrTokenExpired
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.UUID{}, err
	}
	return userID, nil
}
