package ports

import (
	"github.com/google/uuid"
	"github.com/ncfex/dcart/auth-service/internal/domain"
)

type UserRepository interface {
	FindByUsername(username string) (*domain.User, error)
	Create(user *domain.User) error
}

type TokenRepository interface {
	StoreToken(userID uuid.UUID, token string) error
	ValidateToken(token string) (uuid.UUID, error)
}
