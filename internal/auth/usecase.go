package auth

import (
	"context"

	"github.com/kmx0/GophKeeper/internal/models"
)

const CtxUserKey = "user"

type UseCase interface {
	SignUp(ctx context.Context, login, password string) error
	SignIn(ctx context.Context, login, password string) (string, error)
	ParseToken(ctx context.Context, accessToken string) (*models.User, error)
}
