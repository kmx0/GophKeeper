package localstorage

import (
	"context"
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
		return auth.ErrLoginBusy
	}
	s.users[user.Login] = user
	return nil
}

func (s *UserLocalStorage) GetUser(ctx context.Context, login, password string) (*models.User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, user := range s.users {
		if user.Login == login && user.Password == password {
			return user, nil
		}
	}
	return nil, auth.ErrUserNotFound
}
