package users

import (
	"context"
	"fazz/app/internal/domain"
)

type Service interface {
	GetListUsers(ctx context.Context) ([]domain.User, error)
	CreateUser(ctx context.Context, user domain.CreateUser) (domain.User, error)
	AuthorizationUser(ctx context.Context, user domain.LoginUser) (string, error)
}
