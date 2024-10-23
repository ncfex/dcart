package services

import (
	"fmt"
	"time"

	"github.com/ncfex/dcart/auth-service/internal/core/ports"
	"github.com/ncfex/dcart/auth-service/internal/domain"
)

type service struct {
	userRepo  ports.UserRepository
	tokenRepo ports.TokenRepository
}

func NewAuthService(userRepo ports.UserRepository, tokenRepo ports.TokenRepository) ports.AuthService {
	return &service{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

func (s *service) Register(username, password string) error {
	if _, err := s.userRepo.FindByUsername(username); err == nil {
		return domain.ErrUserAlreadyExists
	}

	user := &domain.User{
		Username: username,
		Password: hashPassword(password),
	}

	return s.userRepo.Create(user)
}

func (s *service) Login(username, password string) (string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil || !checkPassword(password, user.Password) {
		return "", domain.ErrInvalidCredentials
	}

	// TODO - use JWT
	token := fmt.Sprintf("%d:%d", user.ID, time.Now().Unix())
	err = s.tokenRepo.StoreToken(user.ID, token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func hashPassword(password string) string {
	// TODO - use bcrypt
	return fmt.Sprintf("hashed_%s", password)
}

func checkPassword(inputPassword, storedPassword string) bool {
	// TODO - use bcrypt
	return inputPassword == storedPassword
}
