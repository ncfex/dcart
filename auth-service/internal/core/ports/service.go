package ports

import "context"

type AuthService interface {
	Register(ctx context.Context, username string, password string) error
	Login(ctx context.Context, username string, password string) (string, error)
}
