package remote

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io/ioutil"

	"github.com/kmx0/GophKeeper/internal/auth"
	"github.com/kmx0/GophKeeper/internal/models"

	"github.com/dgrijalva/jwt-go/v4"
)

var _ auth.UseCase = (*AuthUseCase)(nil)

type AuthClaims struct {
	jwt.StandardClaims
	User *models.User `json:"user"`
}

type AuthUseCase struct {
	userRepo   auth.UserRepository
	singingKey []byte
	hashSalt   string
	tokenFile  string
}

func NewAuthUseCase(userRepo auth.UserRepository, hashSalt string, singingKey []byte, tokenFile string) *AuthUseCase {
	return &AuthUseCase{
		userRepo:   userRepo,
		singingKey: singingKey,
		hashSalt:   hashSalt,
		tokenFile:  tokenFile,
	}
}

func (a *AuthUseCase) SignUp(ctx context.Context, login, password string) error {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(a.hashSalt))
	user := &models.User{
		Login:    login,
		Password: fmt.Sprintf("%x", pwd.Sum(nil)),
	}
	return a.userRepo.CreateUser(ctx, user)
}

// TODO
// SignIn

func (a *AuthUseCase) SignIn(ctx context.Context, login, password string) (string, error) {
	pwd := sha1.New()

	pwd.Write([]byte(password))
	pwd.Write([]byte(a.hashSalt))
	password = fmt.Sprintf("%x", pwd.Sum(nil))

	user, err := a.userRepo.GetUser(ctx, login, password)
	if err != nil {
		return "", auth.ErrUserNotFound
	}
	err = ioutil.WriteFile(a.tokenFile, []byte(user.Token), 0777)
	if err != nil {
		return "", err
	}
	return user.Token, err
}

func (a *AuthUseCase) ParseToken(ctx context.Context, accessToken string) (*models.User, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.singingKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		claims.User.Token = token.Raw
		return claims.User, nil
	}
	return nil, auth.ErrInvalidAccessToken
}
