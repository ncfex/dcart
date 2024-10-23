package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	currentTime := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    "dcart", // move to env etc.
		IssuedAt:  jwt.NewNumericDate(currentTime.UTC()),
		ExpiresAt: jwt.NewNumericDate(currentTime.Add(expiresIn)),
		Subject:   userID.String(),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	// check issuer
	// add more checks
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(t *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return uuid.UUID{}, err
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return uuid.UUID{}, ErrTokenExpired
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

func HashPassword(password string) (string, error) {
	data, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

var (
	ErrTokenExpired = errors.New("token expired")
)
