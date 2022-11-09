package secret

import (
	"context"
	"time"
)

type service struct {
	storage Storage
}

func NewService(storage Storage) Service {
	return &service{storage: storage}
}

func (s *service) GetSecretByLogin(ctx context.Context, owner string) *Secret {
	return s.storage.GetOne(owner)
}

func (s *service) GetAllSecrets(ctx context.Context, updatedAt time.Duration) []*Secret {
	return s.storage.GetAll(updatedAt)
}

func (s *service) CreateSecret(ctx context.Context, dto *CreateSecretDTO) *Secret {
	return nil
}
