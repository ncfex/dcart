package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ncfex/dcart/auth-service/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUserByID(ctx context.Context, userID *uuid.UUID) (*domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
}

type TokenRepository interface {
	StoreToken(ctx context.Context, userID *uuid.UUID, token string) error
	ValidateToken(ctx context.Context, token string) (*uuid.UUID, error)
	RevokeToken(ctx context.Context, token string) error
}
