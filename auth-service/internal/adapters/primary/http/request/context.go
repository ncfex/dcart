package request

import (
	"context"

	"github.com/ncfex/dcart/auth-service/internal/domain"
)

type contextKey string

const UserContextKey contextKey = "user"
const UserIDContextKey contextKey = "userID"

func SetValueToContext(ctx context.Context, key contextKey, value interface{}) context.Context {
	return context.WithValue(ctx, key, value)
}

func GetUserFromContext(ctx context.Context) (*domain.User, bool) {
	user, ok := ctx.Value(UserContextKey).(*domain.User)
	return user, ok
}
