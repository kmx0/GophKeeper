package usecase

import (
	"context"

	"github.com/kmx0/GophKeeper/internal/models"
	"github.com/kmx0/GophKeeper/internal/secret"
	"github.com/stretchr/testify/mock"
)

type SecretUseCaseMock struct {
	mock.Mock
}

var _ secret.UseCase = (*SecretUseCaseMock)(nil)

func (m SecretUseCaseMock) CreateSecret(ctx context.Context, user *models.User, key, value, secretType string) error {
	args := m.Called(user, key, value, secretType)
	return args.Error(0)
}

func (m SecretUseCaseMock) GetSecrets(ctx context.Context, user *models.User) ([]*models.Secret, error) {
	args := m.Called(user)
	return args.Get(0).([]*models.Secret), args.Error(1)
}

func (m SecretUseCaseMock) GetSecret(ctx context.Context, user *models.User, key string) (*models.Secret, error) {
	args := m.Called(user, key)
	return args.Get(0).(*models.Secret), args.Error(1)
}

func (m SecretUseCaseMock) DeleteSecret(ctx context.Context, user *models.User, key string) error {
	args := m.Called(user, key)
	return args.Error(0)
}