package middleware

import (
	"context"
	"net/http"

	"github.com/ncfex/dcart/api-gateway/internal/infrastructure/config"
	"github.com/ncfex/dcart/api-gateway/pkg/httputil/request"
)

type AuthMiddleware struct {
	cfg config.AuthConfig
}

func NewAuthMiddleware(cfg config.AuthConfig) *AuthMiddleware {
	return &AuthMiddleware{
		cfg: cfg,
	}
}

func (am *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), am.cfg.Timeout)
		defer cancel()

		valid, err := am.validateToken(ctx, r.Header)
		if err != nil {
			http.Error(w, "Authorization service unavailable", http.StatusServiceUnavailable)
			return
		}
		if !valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (am *AuthMiddleware) validateToken(ctx context.Context, header http.Header) (bool, error) {
	token, err := request.GetBearerToken(header)
	if err != nil {
		return false, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, am.cfg.ServiceURL, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{
		Timeout: am.cfg.Timeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return false, ctx.Err()
		}
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	return true, nil
}
