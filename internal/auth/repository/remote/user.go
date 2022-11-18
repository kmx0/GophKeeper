package remote

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kmx0/GophKeeper/internal/auth"
	"github.com/kmx0/GophKeeper/internal/models"
	"github.com/sirupsen/logrus"
)

type UserRepository struct {
	client      *http.Client
	repoAddress string
}

var _ auth.UserRepository = (*UserRepository)(nil)

type signInResponse struct {
	Token string `json:"token"`
}

type signInput struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func NewUserRepository(client *http.Client, repoAddress string) *UserRepository {
	return &UserRepository{
		client:      client,
		repoAddress: repoAddress,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	bodyBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}
	bodyIOReader := bytes.NewReader(bodyBytes)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/auth/sign-up", r.repoAddress), bodyIOReader)
	if err != nil {
		logrus.Error(err)
		return err
	}
	request.Header.Add("Content-Type", "application/json")
	// request.Header.Add("X-Real-IP", "application/json")
	response, err := r.client.Do(request)
	if err != nil {
		return err
	}
	// печатаем код ответа
	defer response.Body.Close()
	// читаем поток из тела ответа
	if response.StatusCode > 200 {
		return fmt.Errorf(response.Status)
	}
	return nil

}

func (r *UserRepository) GetUser(ctx context.Context, login, password string) (*models.User, error) {
	sendBody := new(signInput)
	sendBody.Login = login
	sendBody.Password = password
	respToken := new(signInResponse)
	bodyBytes, err := json.Marshal(sendBody)
	if err != nil {
		return nil, err
	}
	bodyIOReader := bytes.NewReader(bodyBytes)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/auth/sign-in", r.repoAddress), bodyIOReader)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	// request.Header.Add("X-Real-IP", "application/json")
	response, err := r.client.Do(request)
	if err != nil {
		return nil, err
	}
	// печатаем код ответа
	defer response.Body.Close()
	// читаем поток из тела ответа
	tokenBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(tokenBody, respToken)
	if err != nil {
		return nil, err
	}

	return &models.User{
		Login:    login,
		Password: password,
		Token:    respToken.Token,
	}, nil

}
