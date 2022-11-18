package usecase

import (
	"context"

	"github.com/kmx0/GophKeeper/internal/auth"
	"github.com/kmx0/GophKeeper/internal/models"
	"github.com/stretchr/testify/mock"
)

type AuthUseCaseMock struct {
	mock.Mock
}

var _ auth.UseCase = (*AuthUseCaseMock)(nil)

func (m *AuthUseCaseMock) SignUp(ctx context.Context, login, password string) error {
	args := m.Called(login, password)
	return args.Error(0)
}

func (m *AuthUseCaseMock) SignIn(ctx context.Context, login, password string) (string, error) {
	args := m.Called(login, password)
	return args.Get(0).(string), args.Error(1)
}

func (m *AuthUseCaseMock) ParseToken(ctx context.Context, accessToken string) (*models.User, error) {
	args := m.Called(accessToken)
	return args.Get(0).(*models.User), args.Error(1)
}
