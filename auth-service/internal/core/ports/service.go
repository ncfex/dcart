package ports

import (
	"context"

	"github.com/ncfex/dcart/auth-service/internal/domain"
)

type AuthService interface {
	Register(ctx context.Context, username string, password string) (*domain.User, error)
	Login(ctx context.Context, username string, password string) (string, error)
}
