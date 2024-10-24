package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	httpAdapter "github.com/ncfex/dcart/auth-service/internal/adapters/primary/http"
	"github.com/ncfex/dcart/auth-service/internal/adapters/primary/http/response"
	"github.com/ncfex/dcart/auth-service/internal/adapters/secondary/postgres"
	"github.com/ncfex/dcart/auth-service/internal/adapters/secondary/redis"
	"github.com/ncfex/dcart/auth-service/internal/core/services/auth"
	"github.com/ncfex/dcart/auth-service/internal/core/services/jsonwebtoken"
	"github.com/ncfex/dcart/auth-service/internal/core/services/password"
	"github.com/ncfex/dcart/auth-service/internal/infrastructure/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	postgresURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB,
	)

	redisURL := fmt.Sprintf(
		"%s:%s",
		cfg.RedisHost,
		cfg.RedisPort,
	)

	// repo
	userRepo, err := postgres.NewUserRepository(postgresURL)
	if err != nil {
		log.Fatalf("Failed to initialize user repository: %v", err)
	}

	tokenRepo, err := redis.NewTokenRepository(redisURL)
	if err != nil {
		log.Fatalf("Failed to initialize token repository: %v", err)
	}

	passwordService := password.NewPasswordService(0)
	jwtService := jsonwebtoken.NewJWTService("dcart", cfg.JwtSecret)
	authService := auth.NewAuthService(userRepo, tokenRepo, passwordService, jwtService)

	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	responder := response.NewHTTPResponder(logger)
	handler := httpAdapter.NewHandler(responder, authService)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler.Router(),
	}

	log.Printf("starting auth service on port %s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
