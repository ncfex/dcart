package authentication

import (
	"context"
	"time"

	"github.com/ncfex/dcart/auth-service/internal/core/ports"
	"github.com/ncfex/dcart/auth-service/internal/domain"
	"github.com/ncfex/dcart/auth-service/internal/domain/errors"
)

type service struct {
	userRepo          ports.UserRepository
	tokenRepo         ports.TokenRepository
	passwordEncrypter ports.PasswordEncrypter
	tokenManager      ports.TokenManager
}

func NewAuthService(
	userRepo ports.UserRepository,
	tokenRepo ports.TokenRepository,
	passwordEncrypter ports.PasswordEncrypter,
	tokenManager ports.TokenManager,
) ports.UserAuthenticator {
	return &service{
		userRepo:          userRepo,
		tokenRepo:         tokenRepo,
		passwordEncrypter: passwordEncrypter,
		tokenManager:      tokenManager,
	}
}

func (s *service) Register(ctx context.Context, username, password string) (*domain.User, error) {
	if username == "" || password == "" {
		return &domain.User{}, errors.ErrInvalidCredentials
	}
	if _, err := s.userRepo.GetUserByUsername(ctx, username); err == nil {
		return &domain.User{}, errors.ErrInvalidCredentials
	}

	hashedPassword, err := s.passwordEncrypter.Hash(password)
	if err != nil {
		return &domain.User{}, domain.ErrUserAlreadyExists
	}

	user := &domain.User{
		Username:     username,
		PasswordHash: hashedPassword,
	}

	return s.userRepo.CreateUser(ctx, user)
}

func (s *service) Login(ctx context.Context, username, password string) (string, error) {
	if username == "" || password == "" {
		return "", domain.ErrInvalidCredentials
	}
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	err = s.passwordEncrypter.Compare(user.PasswordHash, password)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	// TODO - use JWT
	token, err := s.tokenManager.Make(&user.ID, time.Hour*24)
	if err != nil {
		return "", err
	}
	err = s.tokenRepo.StoreToken(ctx, &user.ID, token)
	if err != nil {
		return "", err
	}

	return token, nil
}
