package remote

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kmx0/GophKeeper/internal/models"
	"github.com/kmx0/GophKeeper/internal/secret"
	"github.com/sirupsen/logrus"
)

type SecretRepository struct {
	client      *http.Client
	repoAddress string
}

var _ secret.Repository = (*SecretRepository)(nil)

type Secret struct {
	//int
	ID     string
	UserID string
	Key    string
	Value  string
	Type   string
}
type Secrets struct {
	Secrets []*Secret `json:"secrets"`
}

func NewSecretRepository(client *http.Client, repoAddress string) *SecretRepository {
	return &SecretRepository{
		client:      client,
		repoAddress: repoAddress,
	}
}

func (r *SecretRepository) CreateSecret(ctx context.Context, user *models.User, sc *models.Secret) error {
	bodyBytes, err := json.Marshal(sc)
	if err != nil {
		return err
	}
	bodyIOReader := bytes.NewReader(bodyBytes)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/api/secret/create", r.repoAddress), bodyIOReader)
	if err != nil {
		logrus.Error(err)
		return err
	}
	request.Header.Add("Content-Type", "application/json")
	// request.Header.Add("X-Real-IP", "application/json")

	bearer := "Bearer " + user.Token
	request.Header.Add("Authorization", bearer)
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

type Input struct {
	Key string `json:"key"`
}

func (r *SecretRepository) DeleteSecret(ctx context.Context, user *models.User, key string) error {
	dInput := new(Input)
	dInput.Key = key
	bodyBytes, err := json.Marshal(dInput)
	if err != nil {
		return err
	}
	bodyIOReader := bytes.NewReader(bodyBytes)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/api/secret/delete", r.repoAddress), bodyIOReader)
	if err != nil {
		logrus.Error(err)
		return err
	}
	request.Header.Add("Content-Type", "application/json")
	// request.Header.Add("X-Real-IP", "application/json")

	bearer := "Bearer " + user.Token
	request.Header.Add("Authorization", bearer)
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

type getResponseSingle struct {
	Secret *Secret `json:"secret"`
}

func (r *SecretRepository) GetSecret(ctx context.Context, user *models.User, key string) (*models.Secret, error) {
	getInput := new(Input)
	getInput.Key = key
	bodyBytes, err := json.Marshal(getInput)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	bodyIOReader := bytes.NewReader(bodyBytes)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/api/secret/get", r.repoAddress), bodyIOReader)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	// request.Header.Add("X-Real-IP", "application/json")

	bearer := "Bearer " + user.Token
	request.Header.Add("Authorization", bearer)
	response, err := r.client.Do(request)
	if err != nil {
		return nil, err
	}
	// печатаем код ответа
	defer response.Body.Close()
	// читаем поток из тела ответа
	if response.StatusCode > 200 {
		if response.StatusCode == http.StatusNoContent {
			return nil, secret.ErrSecretNotFound
		}
		return nil, fmt.Errorf("err: %s", response.Status)
	}
	secretBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	respSecret := new(getResponseSingle)
	logrus.Info(string(secretBody))
	err = json.Unmarshal(secretBody, respSecret)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	logrus.Info(respSecret.Secret.Value)
	return &models.Secret{
		ID:     respSecret.Secret.ID,
		UserID: respSecret.Secret.UserID,
		Value:  respSecret.Secret.Value,
		Key:    key,
		Type:   respSecret.Secret.Type,
	}, nil
}

func (r *SecretRepository) GetSecrets(ctx context.Context, user *models.User) ([]*models.Secret, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/api/secret/list", r.repoAddress), nil)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	// request.Header.Add("X-Real-IP", "application/json")

	bearer := "Bearer " + user.Token
	request.Header.Add("Authorization", bearer)
	response, err := r.client.Do(request)
	if err != nil {
		logrus.Error(err)

		return nil, err
	}
	// печатаем код ответа
	defer response.Body.Close()
	// читаем поток из тела ответа
	secretsBody, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Error(err)

		return nil, err
	}
	if response.StatusCode > 200 {
		if response.StatusCode == http.StatusNoContent {
			return nil, secret.ErrUserHaveNotSecret
		}
		return nil, fmt.Errorf("err GetSecrest :%s", response.Status)
	}
	logrus.Info(string(secretsBody))
	respSecrets := new(Secrets)
	err = json.Unmarshal(secretsBody, respSecrets)
	if err != nil {
		logrus.Error(err)

		return nil, err
	}

	return toSecrets(respSecrets), nil
}

func toSecrets(ss *Secrets) []*models.Secret {
	scs := make([]*models.Secret, len(ss.Secrets))
	for i, v := range ss.Secrets {
		scs[i] = &models.Secret{
			ID:     v.ID,
			UserID: v.UserID,
			Key:    v.Key,
			Value:  v.Value,
			Type:   v.Type,
		}
	}
	return scs
}
