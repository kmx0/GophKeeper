package user

import (
	"time"

	"github.com/kmx0/GophKeeper/internal/domain/user"
)

type userStorage struct {
	//client postgres from pkg
}

func (us *userStorage) GetOne(owner string) *user.User {
	return nil
}
func (us *userStorage) GetAll(expiredAt time.Duration) []*user.User {
	return nil
}
func (us *userStorage) Create(user *user.User) *user.User {
	return nil
}
func (us *userStorage) Delete(user *user.User) error {
	return nil
}
