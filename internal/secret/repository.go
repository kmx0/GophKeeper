package secret

import (
	"context"

	"github.com/kmx0/GophKeeper/internal/models"
)

type Repository interface {
	CreateSecret(ctx context.Context, user *models.User, sc *models.Secret) error
	GetSecrets(ctx context.Context, user *models.User) ([]*models.Secret, error)
	GetSecret(ctx context.Context, user *models.User, key string) (*models.Secret, error)
	DeleteSecret(ctx context.Context, user *models.User, key string) error
}
