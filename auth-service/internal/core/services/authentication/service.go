package authentication

import (
	"context"
	"fmt"
	"time"

	"github.com/ncfex/dcart/auth-service/internal/core/ports"
	"github.com/ncfex/dcart/auth-service/internal/core/services/refresh"
	"github.com/ncfex/dcart/auth-service/internal/domain"
	"github.com/ncfex/dcart/auth-service/internal/domain/errors"
)

type service struct {
	userRepo            ports.UserRepository
	tokenRepo           ports.TokenRepository
	passwordEncrypter   ports.PasswordEncrypter
	accessTokenManager  ports.TokenManager
	refreshTokenManager refresh.HexTokenService
}

func NewAuthService(
	userRepo ports.UserRepository,
	tokenRepo ports.TokenRepository,
	passwordEncrypter ports.PasswordEncrypter,
	accessTokenManager ports.TokenManager,
	refreshTokenManager refresh.HexTokenService,
) ports.UserAuthenticator {
	return &service{
		userRepo:            userRepo,
		tokenRepo:           tokenRepo,
		passwordEncrypter:   passwordEncrypter,
		accessTokenManager:  accessTokenManager,
		refreshTokenManager: refreshTokenManager,
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

func (s *service) Login(ctx context.Context, username, password string) (*domain.TokenPair, error) {
	if username == "" || password == "" {
		return nil, domain.ErrInvalidCredentials
	}
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	err = s.passwordEncrypter.Compare(user.PasswordHash, password)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	accessToken, err := s.accessTokenManager.Make(&user.ID, time.Minute*15)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.refreshTokenManager.Make()
	if err != nil {
		return nil, err
	}

	err = s.tokenRepo.StoreToken(ctx, &user.ID, refreshToken)
	if err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  domain.Token(accessToken),
		RefreshToken: domain.Token(refreshToken),
	}, nil
}

func (s *service) Logout(ctx context.Context, token string) error {
	if token == "" {
		return errors.ErrInvalidToken
	}

	err := s.tokenRepo.RevokeToken(ctx, token)
	if err != nil {
		return fmt.Errorf("error revoking token")
	}

	return nil
}
