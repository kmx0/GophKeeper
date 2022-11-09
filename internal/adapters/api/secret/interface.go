package secret

import (
	"context"
	"time"
)

type Service interface {
	GetSecretByLogin(ctx context.Context, owner string) *Secret
	GetAllSecrets(ctx context.Context, updatedAt time.Duration) []*Secret //получить все секреты, старше бла бла
	CreateSecret(ctx context.Context, dto *CreateSecretDTO) *Secret
}
