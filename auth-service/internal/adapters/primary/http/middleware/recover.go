package middleware

import (
	"errors"
	"net/http"

	"github.com/ncfex/dcart/auth-service/internal/adapters/primary/http/response"
)

var ErrInternalServerErrorStr = "internal server error"
var ErrInternalServerError = errors.New(ErrInternalServerErrorStr)

func Recovery(responder response.Responder) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					responder.RespondWithError(
						w,
						http.StatusInternalServerError,
						ErrInternalServerErrorStr,
						ErrInternalServerError,
					)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
