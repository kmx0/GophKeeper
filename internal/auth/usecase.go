package auth

import (
	"context"

	"github.com/kmx0/GophKeeper/internal/models"
)

type UseCase interface{
	SignUp(ctx context.Context, username, password string)error
	SignIn(ctx context.Context, username, password string) (string, error)
	ParseToken(ctx context.Context, accesaToken string)(*models.User, error)
}