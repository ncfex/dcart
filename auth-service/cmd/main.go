package main

import (
	"fmt"
	"log"
	"net/http"

	httpAdapter "github.com/ncfex/dcart/auth-service/internal/adapters/primary/http"
	"github.com/ncfex/dcart/auth-service/internal/adapters/secondary/postgres"
	"github.com/ncfex/dcart/auth-service/internal/adapters/secondary/redis"
	"github.com/ncfex/dcart/auth-service/internal/core/services"
	"github.com/ncfex/dcart/auth-service/internal/infrastructure/config"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Println(cfg)

	// repo
	userRepo := postgres.NewUserRepository(cfg.PostgresURL)
	tokenRepo := redis.NewTokenRepository(cfg.RedisURL)
	authService := services.NewAuthService(userRepo, tokenRepo)

	handler := httpAdapter.NewHandler(authService)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler.Router(),
	}

	log.Printf("starting auth service on port %s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
