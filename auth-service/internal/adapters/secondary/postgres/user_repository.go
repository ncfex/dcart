package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/ncfex/dcart/auth-service/internal/core/ports"
	"github.com/ncfex/dcart/auth-service/internal/domain"
	database "github.com/ncfex/dcart/auth-service/internal/infrastructure/database/sqlc"

	_ "github.com/lib/pq"
)

type repository struct {
	db *database.Queries
}

func NewUserRepository(dsn string) (ports.UserRepository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error initializing database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	queries := database.New(db)
	return &repository{
		db: queries,
	}, nil
}

func (r *repository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	dbUser, err := r.db.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return domain.NewUserFromDB(&dbUser), nil
}

func (r *repository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	params := database.CreateUserParams{
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
	}

	dbUser, err := r.db.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}
	return domain.NewUserFromDB(&dbUser), nil
}

func (r *repository) GetUserByID(ctx context.Context, userID *uuid.UUID) (*domain.User, error) {
	dbUser, err := r.db.GetUserByID(ctx, *userID)
	if err != nil {
		return nil, err
	}

	return domain.NewUserFromDB(&dbUser), nil
}
