package usecase

import (
	"context"
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/kmx0/GophKeeper/internal/auth"
	"github.com/kmx0/GophKeeper/internal/models"
)
type AuthClaims struct {
	jwt.StandartClaims
	User *models.User `json:"user"`
}

type AuthUseCase struct {
	userRepo       auth.UserRepository
	hashSalt       string
	singingKey     []byte
	expireDuration time.Duration
}

func NewAuthUseCase(
	userRepo auth.UserRepository,
	hashSalt string,
	singingKey []byte,
	tokenTTL time.Duration) *AuthUseCase {
	return &AuthUseCase{
		userRepo:       userRepo,
		hashSalt:       hashSalt,
		singingKey:     singingKey,
		expireDuration: time.Second * tokenTTL,
	}
}

func (a *AuthUseCase) SignUp(ctx context.Context, login, password string) error {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(a.hashSalt))
	user := &models.User{
		Login: login,
		Password: fmt.Sprintf("%x", pwd.Sum(nil)),
	}
	return a.userRepo.CreateUser(ctx, user)
}
// TODO
 //SignIn

 //ParseToken