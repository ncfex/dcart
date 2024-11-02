package middleware

import (
	"net/http"

	"github.com/ncfex/dcart/auth-service/internal/adapters/primary/http/request"
	"github.com/ncfex/dcart/auth-service/internal/core/ports"
)

func AuthenticateWithJWT(
	tokenManager ports.TokenManager,
	tokenRepo ports.TokenRepository,
	userRepo ports.UserRepository,
) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

			user, err := userRepo.GetUserByID(r.Context(), userID)
			if err != nil {
				http.Error(w, "not authorized", http.StatusUnauthorized)
				return
			}

			ctx := request.SetValueToContext(r.Context(), request.UserIDContextKey, userID)
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
			refreshToken, err := request.GetBearerToken(r.Header)
			if err != nil {
				http.Error(w, "not authorized", http.StatusUnauthorized)
				return
			}

			user, err := tokenRepo.GetUserFromToken(r.Context(), refreshToken)
			if err != nil {
				http.Error(w, "not authorized", http.StatusUnauthorized)
				return
			}

			ctx := request.SetValueToContext(r.Context(), request.UserIDContextKey, user.ID)
			ctx = request.SetValueToContext(ctx, request.UserContextKey, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
