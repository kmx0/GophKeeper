package usecase

import (
	"context"
	"crypto/sha1"
	"fmt"
	"strconv"
	"time"

	"github.com/kmx0/GophKeeper/internal/auth"
	"github.com/kmx0/GophKeeper/internal/models"

	"github.com/dgrijalva/jwt-go/v4"
)

var _ auth.UseCase = (*AuthUseCase)(nil)

type AuthClaims struct {
	jwt.StandardClaims
	User *User `json:"user"`
}

type User struct {
	//int
	ID       string `json:"id,omitempty"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
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
		Login:    login,
		Password: fmt.Sprintf("%x", pwd.Sum(nil)),
	}
	return a.userRepo.CreateUser(ctx, user)
}

func (a *AuthUseCase) SignIn(ctx context.Context, login, password string) (string, error) {
	pwd := sha1.New()

	pwd.Write([]byte(password))
	pwd.Write([]byte(a.hashSalt))
	password = fmt.Sprintf("%x", pwd.Sum(nil))

	user, err := a.userRepo.GetUser(ctx, login, password)
	if err != nil {
		return "", fmt.Errorf("error on SignIn: %w", auth.ErrUserNotFound)
	}
	claims := AuthClaims{
		User: toUser(user),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(a.expireDuration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.singingKey)
}

func (a *AuthUseCase) ParseToken(ctx context.Context, accessToken string) (*models.User, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error on ParseToken: %w", fmt.Errorf("errunexpected signing method: %v", token.Header["alg"]))
		}
		return a.singingKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("error on ParseToken: %w", err)
	}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return toModelUser(claims.User), nil
	}
	return nil, fmt.Errorf("error on ParseToken: %w", auth.ErrInvalidAccessToken)
}

func toModelUser(user *User) *models.User {
	id, _ := strconv.Atoi(user.ID)
	return &models.User{
		ID:       id,
		Login:    user.Login,
		Password: user.Password,
		Token:    user.Token,
	}
}

func toUser(user *models.User) *User {
	return &User{
		ID:       strconv.Itoa(user.ID),
		Login:    user.Login,
		Password: user.Password,
		Token:    user.Token,
	}
}
