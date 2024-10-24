package jsonwebtoken

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrTokenExpired  = errors.New("token expired")
	ErrInvalidToken  = errors.New("invalid token")
	ErrInvalidIssuer = errors.New("invalid issuer")
	ErrInvalidUserID = errors.New("invalid user ID")
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
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(t *jwt.Token) (interface{}, error) { return []byte(s.tokenSecret), nil },
	)
	if err != nil {
		return uuid.Nil, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}

	if issuer != string(s.issuer) {
		return uuid.Nil, ErrInvalidIssuer
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, ErrInvalidUserID
	}
	return userID, nil
}
