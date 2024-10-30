package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ncfex/dcart/auth-service/internal/domain"
)

type UserAuthenticator interface {
	Register(ctx context.Context, username string, password string) (*domain.User, error)
	Login(ctx context.Context, username string, password string) (string, error)
}

type PasswordEncrypter interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
}

type TokenManager interface {
	Make(userID *uuid.UUID, expiresIn time.Duration) (string, error)
	Validate(token string) (*uuid.UUID, error)
}
