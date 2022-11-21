package cli

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/kmx0/GophKeeper/internal/auth"
	"github.com/kmx0/GophKeeper/internal/models"
)

type AuthMiddleware struct {
	usecase   auth.UseCase
	tokenFile string
}

func NewAuthMiddleware(usecase auth.UseCase, tokenFile string) *AuthMiddleware {
	return (&AuthMiddleware{
		usecase:   usecase,
		tokenFile: tokenFile,
	})
}

func (m *AuthMiddleware) Handle(ctx context.Context) (*models.User, error) {

	tokenFile, err := os.Open(m.tokenFile)
	if err != nil {
		// fmt.Println(err)
		if strings.Contains(err.Error(), "no such file or directory") {
			return nil, auth.ErrUserNotLoggedIn
		}
		return nil, err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer tokenFile.Close()

	byteValue, err := ioutil.ReadAll(tokenFile)
	if err != nil {
		return nil, err
	}
	user, err := m.usecase.ParseToken(ctx, string(byteValue))
	switch err {

	case nil:
		// user.Token =
		return user, nil
	case auth.ErrInvalidAccessToken:
		return nil, auth.ErrInvalidAccessToken
	default:
		if strings.Contains(err.Error(), "token is expired"){
			return nil, fmt.Errorf("please login again: %s", err.Error())
		}
		return nil, err
	}

}
