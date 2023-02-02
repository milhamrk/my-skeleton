package users

import (
	"context"
	"fazz/app/internal/domain"
)

type Database interface {
	GetUsers(ctx context.Context) ([]domain.User, error)
	CreateUser(ctx context.Context, user domain.User, password string) (domain.User, error)
	GetIDAndPasswordByEmail(ctx context.Context, email string) (string, string, error)
}
