package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ncfex/dcart/auth-service/internal/adapters/primary/http/request"
	"github.com/ncfex/dcart/auth-service/internal/core/ports"
)

// TODO consider response.Responder
func AuthenticateWithJWT(
	tokenManager ports.TokenManager,
	tokenRepo ports.TokenRepository,
	userRepo ports.UserRepository,
) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
			defer cancel()

			accessToken, err := request.GetBearerToken(r.Header)
			if err != nil {
				http.Error(w, "not authorized", http.StatusUnauthorized)
				return
			}

			userID, err := tokenManager.Validate(accessToken)
			if err != nil {
				http.Error(w, "not authorized", http.StatusUnauthorized)
				return
			}

			user, err := userRepo.GetUserByID(ctx, userID)
			if err != nil {
				switch {
				case errors.Is(err, context.DeadlineExceeded):
					http.Error(w, "request timeout", http.StatusGatewayTimeout)
				default:
					http.Error(w, "not authorized", http.StatusUnauthorized)
				}
				return
			}

			ctx = request.SetValueToContext(ctx, request.UserIDContextKey, userID)
			ctx = request.SetValueToContext(ctx, request.UserContextKey, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AuthenticateWithRefreshToken(
	tokenManager ports.TokenManager,
	tokenRepo ports.TokenRepository,
	userRepo ports.UserRepository,
) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
			defer cancel()

			refreshToken, err := request.GetBearerToken(r.Header)
			if err != nil {
				http.Error(w, "not authorized", http.StatusUnauthorized)
				return
			}

			user, err := tokenRepo.GetUserFromToken(ctx, refreshToken)
			if err != nil {
				switch {
				case errors.Is(err, context.DeadlineExceeded):
					http.Error(w, "request timeout", http.StatusGatewayTimeout)
				default:
					http.Error(w, "not authorized", http.StatusUnauthorized)
				}
				return
			}

			ctx = request.SetValueToContext(ctx, request.UserIDContextKey, user.ID)
			ctx = request.SetValueToContext(ctx, request.UserContextKey, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
