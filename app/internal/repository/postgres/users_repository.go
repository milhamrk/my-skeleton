package postgres

import (
	"context"
	"fazz/app/internal/domain"

	"github.com/oklog/ulid/v2"
	"github.com/uptrace/bun"
)

type UsersRepository struct {
	db *bun.DB
}

func NewUsersRepository(db *bun.DB) *UsersRepository {
	return &UsersRepository{
		db: db,
	}
}

func (u *UsersRepository) GetUsers(ctx context.Context) ([]domain.User, error) {
	usersModel := make([]userModel, 0)
	err := u.db.NewSelect().Model(&usersModel).Scan(ctx)
	if err != nil {
		return nil, err
	}
	users := userModelsToDomain(usersModel)
	return users, nil
}

func (u *UsersRepository) CreateUser(ctx context.Context, user domain.User, password string) (domain.User, error) {
	model := userModel{}
	model.FromDomain(user)
	model.ID = ulid.Make().String()
	model.EncryptedPassword = password
	_, err := u.db.NewInsert().Model(&model).Exec(ctx)
	if err != nil {
		return userModelToDomain(userModel{}), err
	}
	return userModelToDomain(model), nil
}

func (u *UsersRepository) GetIDAndPasswordByEmail(ctx context.Context, email string) (string, string, error) {
	var password string
	var id string

	model := userModel{}
	err := u.db.NewSelect().Model(&model).
		Where("email = ?", email).Column("id", "encrypted_password").Scan(ctx, &id, &password)
	if err != nil {
		return "", "", err
	}
	return id, password, nil
}
