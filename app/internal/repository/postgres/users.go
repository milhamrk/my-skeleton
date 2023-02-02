package postgres

import (
	"fazz/app/internal/domain"

	"github.com/uptrace/bun"
)

type userModel struct {
	bun.BaseModel `bun:"table:users,alias:users,select:users"`

	ID                string `bun:"id,pk"`
	Email             string `bun:"email,unique"`
	Name              string `bun:"name"`
	Surname           string `bun:"surname"`
	Patronymic        string `bun:"patronymic"`
	Role              string `json:"role"`
	EncryptedPassword string `bun:"encrypted_password"`
}

func (u *userModel) FromDomain(user domain.User) {
	u.ID = user.ID
	u.Email = user.Email
	u.Name = user.Name
	u.Surname = user.Surname
	u.Patronymic = user.Patronymic
	u.Role = user.Role
}

func userModelToDomain(model userModel) domain.User {
	return domain.User{
		ID:         model.ID,
		Email:      model.Email,
		Name:       model.Name,
		Surname:    model.Surname,
		Patronymic: model.Patronymic,
		Role:       model.Role,
	}
}

func userModelsToDomain(models []userModel) []domain.User {
	users := make([]domain.User, 0, len(models))
	for _, model := range models {
		users = append(users, userModelToDomain(model))
	}
	return users
}
