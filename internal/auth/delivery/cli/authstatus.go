package cli

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/kmx0/GophKeeper/internal/auth"
	"github.com/kmx0/GophKeeper/internal/models"
)

type AuthStatus struct {
	usecase   auth.UseCase
	tokenFile string
}

func NewAuthStatus(usecase auth.UseCase, tokenFile string) *AuthStatus {
	return (&AuthStatus{
		usecase:   usecase,
		tokenFile: tokenFile,
	})
}

func (m *AuthStatus) CheckAuthStatus(ctx context.Context) (*models.User, error) {

	tokenFile, err := os.Open(m.tokenFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	if errors.Is(err, os.ErrNotExist) {
		return nil, auth.ErrUserNotLoggedIn
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
		if strings.Contains(err.Error(), "token is expired") {
			return nil, fmt.Errorf("please login again: %s", err.Error())
		}
		return nil, err
	}

}
