package services

import (
	"context"
	"fmt"
	"time"

	"github.com/ncfex/dcart/auth-service/internal/core/ports"
	"github.com/ncfex/dcart/auth-service/internal/core/services/auth"
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

func (s *service) Register(ctx context.Context, username, password string) error {
	if username == "" || password == "" {
		return domain.ErrInvalidCredentials
	}
	if _, err := s.userRepo.FindByUsername(username); err == nil {
		return domain.ErrUserAlreadyExists
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return domain.ErrUserAlreadyExists
	}

	user := &domain.User{
		Username:     username,
		PasswordHash: hashedPassword,
	}

	return s.userRepo.Create(user)
}

func (s *service) Login(ctx context.Context, username, password string) (string, error) {
	if username == "" || password == "" {
		return "", domain.ErrInvalidCredentials
	}
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	err = auth.CheckPasswordHash(password, user.PasswordHash)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	// TODO - use JWT
	token := fmt.Sprintf("%s:%d", user.ID.String(), time.Now().Unix())
	err = s.tokenRepo.StoreToken(user.ID, token)
	if err != nil {
		return "", err
	}

	return token, nil
}
