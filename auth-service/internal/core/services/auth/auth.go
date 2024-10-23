package auth

import (
	"context"
	"time"

	"github.com/ncfex/dcart/auth-service/internal/core/ports"
	"github.com/ncfex/dcart/auth-service/internal/domain"
	"github.com/ncfex/dcart/auth-service/internal/domain/errors"
)

type service struct {
	userRepo        ports.UserRepository
	tokenRepo       ports.TokenRepository
	passwordService *PasswordService
	jwtService      *JWTService
}

func NewAuthService(
	userRepo ports.UserRepository,
	tokenRepo ports.TokenRepository,
	passwordService *PasswordService,
	jwtService *JWTService,
) ports.AuthService {
	return &service{
		userRepo:        userRepo,
		tokenRepo:       tokenRepo,
		passwordService: passwordService,
		jwtService:      jwtService,
	}
}

func (s *service) Register(ctx context.Context, username, password string) (*domain.User, error) {
	if username == "" || password == "" {
		return &domain.User{}, errors.ErrInvalidCredentials
	}
	if _, err := s.userRepo.FindByUsername(username); err == nil {
		return &domain.User{}, errors.ErrInvalidCredentials
	}

	hashedPassword, err := s.passwordService.HashPassword(password)
	if err != nil {
		return &domain.User{}, domain.ErrUserAlreadyExists
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

	err = s.passwordService.CheckPasswordHash(password, user.PasswordHash)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	// TODO - use JWT
	token, err := s.jwtService.MakeJWT(user.ID, time.Hour*24)
	if err != nil {
		return "", err
	}
	err = s.tokenRepo.StoreToken(user.ID, token)
	if err != nil {
		return "", err
	}

	return token, nil
}
