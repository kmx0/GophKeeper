package user

import (
	"context"
	"time"

	"github.com/kmx0/GophKeeper/internal/domain/user"
)

type Service interface {
	GetByLogin(ctx context.Context, owner string) *user.User
	GetAll(ctx context.Context, updatedAt time.Duration) []*user.User //получить все секреты, старше бла бла
	Create(ctx context.Context, dto *user.CreateUserDTO) *user.User
}
