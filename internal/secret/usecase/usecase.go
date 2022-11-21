package usecase

import (
	"context"

	"github.com/kmx0/GophKeeper/internal/models"
	"github.com/kmx0/GophKeeper/internal/secret"
)

var _ secret.UseCase = (*SecretUseCase)(nil)

type SecretUseCase struct {
	secretRepo secret.Repository
}

func NewSecretUseCase(secretRepo secret.Repository) *SecretUseCase {
	return &SecretUseCase{
		secretRepo: secretRepo,
	}
}
func (s *SecretUseCase) CreateSecret(ctx context.Context, user *models.User, key, value, secretType string) error {
	sc := &models.Secret{
		Key:    key,
		Value:  value,
		UserID: user.ID,
		Type:   secretType,
	}
	return s.secretRepo.CreateSecret(ctx, user, sc)
}

func (s *SecretUseCase) GetSecrets(ctx context.Context, user *models.User) ([]*models.Secret, error) {
	return s.secretRepo.GetSecrets(ctx, user)
}

func (s *SecretUseCase) GetSecret(ctx context.Context, user *models.User, key string) (*models.Secret, error) {
	return s.secretRepo.GetSecret(ctx, user, key)
}
func (s *SecretUseCase) DeleteSecret(ctx context.Context, user *models.User, id string) error {
	return s.secretRepo.DeleteSecret(ctx, user, id)
}
