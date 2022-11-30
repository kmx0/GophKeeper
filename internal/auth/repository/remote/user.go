package remote

import (
	"context"

	"github.com/kmx0/GophKeeper/internal/auth"
	"github.com/kmx0/GophKeeper/internal/models"
)

type UserRepository struct {
	requests auth.Requests
}

var _ auth.UserRepository = (*UserRepository)(nil)

func NewUserRepository(requests auth.Requests) *UserRepository {
	return &UserRepository{
		requests: requests,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	return r.requests.CreateUser(ctx, user)

}

func (r *UserRepository) GetUser(ctx context.Context, login, password string) (*models.User, error) {
	return r.requests.GetUser(ctx, login, password)
}
