package auth

import (
	"context"

	"github.com/kmx0/GophKeeper/internal/models"
)

type Requests interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, login, password string) (*models.User, error)
}
