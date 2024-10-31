package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ncfex/dcart/auth-service/internal/core/ports"
	"github.com/ncfex/dcart/auth-service/internal/infrastructure/database/postgres"
	database "github.com/ncfex/dcart/auth-service/internal/infrastructure/database/sqlc"
)

var (
	ErrTokenNotFound   = errors.New("token not found")
	ErrTokenExpired    = errors.New("token expired")
	ErrTokenRevoked    = errors.New("token revoked")
	ErrInvalidToken    = errors.New("invalid token")
	ErrStoringToken    = errors.New("error storing token")
	ErrValidatingToken = errors.New("error validating token")
)

type tokenRepository struct {
	queries   *database.Queries
	expiresIn time.Duration
}

func NewTokenRepository(db *postgres.Database, expiresIn time.Duration) ports.TokenRepository {
	return &tokenRepository{
		queries:   database.New(db.DB),
		expiresIn: expiresIn,
	}
}

func (r *tokenRepository) StoreToken(ctx context.Context, userID *uuid.UUID, token string) error {
	params := database.CreateRefreshTokenParams{
		Token:     token,
		UserID:    *userID,
		ExpiresAt: time.Now().Add(r.expiresIn),
	}

	_, err := r.queries.CreateRefreshToken(ctx, params)
	if err != nil {
		return errors.Join(ErrStoringToken, err)
	}

	return nil
}

func (r *tokenRepository) ValidateToken(ctx context.Context, token string) (*uuid.UUID, error) {
	user, err := r.queries.GetUserFromRefreshToken(ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTokenNotFound
		}
		return nil, errors.Join(ErrValidatingToken, err)
	}

	return &user.ID, nil
}

func (r *tokenRepository) RevokeToken(ctx context.Context, token string) error {
	_, err := r.queries.RevokeRefreshToken(ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrTokenNotFound
		}
		return err
	}

	return nil
}
