package localstorage

import (
	"context"
	"fmt"
	"sync"

	"github.com/kmx0/GophKeeper/internal/auth"
	"github.com/kmx0/GophKeeper/internal/models"
)

type UserLocalStorage struct {
	users map[string]*models.User
	mutex *sync.Mutex
}

var _ auth.UserRepository = (*UserLocalStorage)(nil)

func NewUserLocalStorage() *UserLocalStorage {
	return &UserLocalStorage{
		users: make(map[string]*models.User),
		mutex: new(sync.Mutex),
	}
}

func (s *UserLocalStorage) CreateUser(ctx context.Context, user *models.User) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, ok := s.users[user.Login]; ok {
		return fmt.Errorf("error on CreateUser: %w", auth.ErrLoginBusy)
	}
	s.users[user.Login] = user
	return nil
}

func (s *UserLocalStorage) GetUser(ctx context.Context, login, password string) (*models.User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	user := s.users[login]
	if user == nil {
		return nil, fmt.Errorf("error on GetUser: %w", auth.ErrUserNotFound)
	}
	return user, nil
}
