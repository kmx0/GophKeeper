package secret

import (
	"context"

	"github.com/kmx0/GophKeeper/internal/models"
)

type UseCase interface {
	CreateSecret(ctx context.Context, user *models.User, key, value, secretType string) error
	GetSecrets(ctx context.Context, user *models.User) ([]*models.Secret, error)
	GetSecret(ctx context.Context, user *models.User, key string) (*models.Secret, error)
	DeleteSecret(ctx context.Context, user *models.User, key string) error
}
