package secret

import "time"

type Storage interface {
	GetOne(owner string) *Secret
	GetAll(expiredAt time.Duration) []*Secret
	Create(secret *Secret) *Secret
	Delete(secret *Secret) error
}
