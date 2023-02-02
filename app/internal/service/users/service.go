package users

import (
	"context"
	"errors"
	"fazz/app/internal/domain"
	"fazz/app/pkg/auth"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	db           Database
	log          *zap.Logger
	tokenManager *auth.Manager
}

// TODO: tokenTTL to config
const tokenTTL = 48 * time.Hour

func NewService(db Database, log *zap.Logger, tokenManager *auth.Manager) *Service {
	return &Service{db: db, log: log, tokenManager: tokenManager}
}

func (s *Service) GetListUsers(ctx context.Context) ([]domain.User, error) {
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		s.log.Error("error get all users", zap.Error(err))
		return nil, err
	}
	return users, nil
}

func (s *Service) CreateUser(ctx context.Context, createUser domain.CreateUser) (domain.User, error) {
	user := domain.User{}
	encryptedPassword, err := domain.EncryptPassword(createUser.Password)
	if err != nil {
		return user, err
	}
	user.FromCreateUser(createUser)
	user, err = s.db.CreateUser(ctx, user, encryptedPassword)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *Service) AuthorizationUser(ctx context.Context, loginUser domain.LoginUser) (string, error) {
	ID, encryptedPassword, err := s.db.GetIDAndPasswordByEmail(ctx, loginUser.Email)
	if err != nil {
		return "", err
	}
	if !domain.ComparePassword(loginUser.Password, encryptedPassword) {
		return "", errors.New("invalid email or password")
	}
	token, err := s.tokenManager.NewJWT(ID, 48*time.Hour)
	if err != nil {
		return "", err
	}
	return token, nil
}
