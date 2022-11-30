package remote

import (
	"context"

	"github.com/kmx0/GophKeeper/internal/models"
	"github.com/kmx0/GophKeeper/internal/secret"
)

type SecretRepository struct {
	requests secret.Requests
}

var _ secret.Repository = (*SecretRepository)(nil)

func NewSecretRepository(requests secret.Requests) *SecretRepository {
	return &SecretRepository{
		requests: requests,
	}
}

func (r *SecretRepository) CreateSecret(ctx context.Context, user *models.User, sc *models.Secret) error {
	return r.requests.CreateSecret(ctx, user, sc)
}

func (r *SecretRepository) DeleteSecret(ctx context.Context, user *models.User, key string) error {
	return r.requests.DeleteSecret(ctx, user, key)
}

func (r *SecretRepository) GetSecret(ctx context.Context, user *models.User, key string) (*models.Secret, error) {
	return r.requests.GetSecret(ctx, user, key)
}

func (r *SecretRepository) GetSecrets(ctx context.Context, user *models.User) ([]*models.Secret, error) {
	return r.requests.GetSecrets(ctx, user)
}
