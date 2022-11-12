package secret

import (
	"context"

	"github.com/kmx0/GophKeeper/internal/models"
)

type UseCase interface {
	CreateSecret(ctx context.Context, user *models.User, value string) error
	GetSecrets(ctx context.Context, user *models.User) ([]*models.Secret, error)
	GetSecret(ctx context.Context, user *models.User, id string) (*models.Secret, error)
	DeleteSecret(ctx context.Context, user *models.User, id string) error
}
