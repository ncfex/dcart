package password

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmptyPassword = errors.New("password cannot be empty")
)

type PasswordService struct {
	cost int
}

func NewPasswordService(cost int) *PasswordService {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	return &PasswordService{cost: cost}
}

func (s *PasswordService) HashPassword(password string) (string, error) {
	if password == "" {
		return "", ErrEmptyPassword
	}

	data, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *PasswordService) CheckPasswordHash(password, hash string) error {
	if password == "" {
		return ErrEmptyPassword
	}

	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
