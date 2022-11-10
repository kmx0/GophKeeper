package localstorage

import (
	"sync"

	"github.com/kmx0/GophKeeper/internal/models"
)

type UserLocalStorage struct {
	users map[string]*models.User
	mutex *sync.Mutex
}

func NewUserLocalStorage() *UserLocalStorage {
	return &UserLocalStorage{
		users: make(map[string]*models.User),
		mutex: new(sync.Mutex),
	}
}
