package redis

import (
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/ncfex/dcart/auth-service/internal/core/ports"
)

type repository struct {
	pool *redis.Pool
}

func NewTokenRepository(redisURL string) (ports.TokenRepository, error) {
	pool := &redis.Pool{
		MaxIdle:   10,
		MaxActive: 100,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisURL)
		},
		IdleTimeout: 240 * time.Second,
	}

	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("PING")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &repository{pool: pool}, nil
}

func (r *repository) StoreToken(userID uuid.UUID, token string) error {
	conn := r.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SETEX", token, 86400, userID.String()) // ttl 24h
	return err
}

func (r *repository) ValidateToken(token string) (uuid.UUID, error) {
	conn := r.pool.Get()
	defer conn.Close()

	userIDStr, err := redis.String(conn.Do("GET", token))
	if errors.Is(err, redis.ErrNil) {
		return uuid.UUID{}, nil // token not found
	}
	if err != nil {
		return uuid.UUID{}, err
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("invalid UUID format: %v", err)
	}

	return userID, nil
}
