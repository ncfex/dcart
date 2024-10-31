package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/ncfex/dcart/auth-service/internal/core/ports"
	"github.com/ncfex/dcart/auth-service/internal/domain"
	"github.com/ncfex/dcart/auth-service/internal/infrastructure/database/postgres"
	database "github.com/ncfex/dcart/auth-service/internal/infrastructure/database/sqlc"
)

type userRepository struct {
	queries *database.Queries
}

func NewUserRepository(db *postgres.Database) ports.UserRepository {
	return &userRepository{
		queries: database.New(db.DB),
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	params := database.CreateUserParams{
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
	}

	dbUser, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}
	return domain.NewUserFromDB(&dbUser), nil
}

func (r *userRepository) GetUserByID(ctx context.Context, userID *uuid.UUID) (*domain.User, error) {
	dbUser, err := r.queries.GetUserByID(ctx, *userID)
	if err != nil {
		return nil, err
	}
	return domain.NewUserFromDB(&dbUser), nil
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	dbUser, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return domain.NewUserFromDB(&dbUser), nil
}
