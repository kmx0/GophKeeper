package secret

import (
	"time"

	"gophkee.per/internal/domain/secret"
)

type secretStorage struct {
	//client postgres from pkg
}

func (ss *secretStorage) GetOne(owner string) *secret.Secret {
	return nil
}
func (ss *secretStorage) GetAll(expiredAt time.Duration) []*secret.Secret {
	return nil
}
func (ss *secretStorage) Create(secret *secret.Secret) *secret.Secret {
	return nil
}
func (ss *secretStorage) Delete(secret *secret.Secret) error {
	return nil
}
