package localstorage

import (
	"context"
	"strconv"
	"sync"

	"github.com/kmx0/GophKeeper/internal/models"
	"github.com/kmx0/GophKeeper/internal/secret"
)

type SecretLocalStorage struct {
	secrets  map[string]*models.Secret
	secretID int
	mutex    *sync.Mutex
}

var _ secret.Repository = (*SecretLocalStorage)(nil)

func NewSecretLocalStorage() *SecretLocalStorage {
	return &SecretLocalStorage{
		secrets: make(map[string]*models.Secret),
		mutex:   new(sync.Mutex),
	}
}

func (s *SecretLocalStorage) CreateSecret(ctx context.Context, user *models.User, sc *models.Secret) error {
	s.mutex.Lock()
	sc.UserID = user.ID
	s.secrets[strconv.Itoa(s.secretID)] = sc
	s.secretID += 1
	s.mutex.Unlock()
	return nil
}

func (s *SecretLocalStorage) GetSecrets(ctx context.Context, user *models.User) ([]*models.Secret, error) {
	secrets := make([]*models.Secret, 0)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, secret := range s.secrets {
		if secret.UserID == user.ID {
			secrets = append(secrets, secret)
		}
	}
	if len(secrets) == 0 {

		return nil, secret.ErrSecretNotFound
	}
	return secrets, nil

}

func (s *SecretLocalStorage) GetSecret(ctx context.Context, user *models.User, id string) (*models.Secret, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// for _, secret := range s.secrets {
	if secret, ok := s.secrets[id]; ok && secret.UserID == user.ID {
		return secret, nil
	}
	return nil, secret.ErrSecretNotFound
}
func (s *SecretLocalStorage) DeleteSecret(ctx context.Context, user *models.User, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if secret, ok := s.secrets[id]; ok && secret.UserID == user.ID {
		delete(s.secrets, id)
	}
	return secret.ErrSecretNotFound
}
