package mock

import (
	"context"

	"github.com/kmx0/GophKeeper/internal/auth"
	"github.com/kmx0/GophKeeper/internal/models"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

var _ auth.UserRepository = (*UserRepositoryMock)(nil)

func (s *UserRepositoryMock) CreateUser(ctx context.Context, user *models.User) error {
	args := s.Called(user)

	return args.Error(0)
}

func (s *UserRepositoryMock) GetUser(ctx context.Context, login, password string) (*models.User, error) {
	args := s.Called(login, password)

	return args.Get(0).(*models.User), args.Error(1)
}
