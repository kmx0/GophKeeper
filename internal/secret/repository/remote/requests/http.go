package requests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kmx0/GophKeeper/internal/models"
	"github.com/kmx0/GophKeeper/internal/secret"
)

type SecretRequests struct {
	client      *http.Client
	repoAddress string
}

var _ secret.Requests = (*SecretRequests)(nil)

type Secret struct {
	//int
	ID     int
	UserID int
	Key    string
	Value  string
	Type   string
}
type Secrets struct {
	Secrets []*Secret `json:"secrets"`
}

func NewSecretRequest(client *http.Client, repoAddress string) *SecretRequests {
	return &SecretRequests{
		client:      client,
		repoAddress: repoAddress,
	}
}

func (r *SecretRequests) CreateSecret(ctx context.Context, user *models.User, sc *models.Secret) error {
	bodyBytes, err := json.Marshal(sc)
	if err != nil {
		return fmt.Errorf("err on CreateSecret: %w", err)
	}
	bodyIOReader := bytes.NewReader(bodyBytes)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/api/secret/create", r.repoAddress), bodyIOReader)
	if err != nil {
		return fmt.Errorf("err on CreateSecret: %w", err)
	}
	request.Header.Add("Content-Type", "application/json")
	// request.Header.Add("X-Real-IP", "application/json")

	bearer := "Bearer " + user.Token
	request.Header.Add("Authorization", bearer)
	response, err := r.client.Do(request)
	if err != nil {
		return fmt.Errorf("err on CreateSecret: %w", err)
	}
	// печатаем код ответа
	defer response.Body.Close()
	// читаем поток из тела ответа
	if response.StatusCode > 200 {
		return fmt.Errorf("err on CreateSecret: %w", fmt.Errorf(response.Status))
	}
	return nil

}

type Input struct {
	Key string `json:"key"`
}

func (r *SecretRequests) DeleteSecret(ctx context.Context, user *models.User, key string) error {
	dInput := new(Input)
	dInput.Key = key
	bodyBytes, err := json.Marshal(dInput)
	if err != nil {
		return fmt.Errorf("err on DeleteSecret :%w", err)
	}
	bodyIOReader := bytes.NewReader(bodyBytes)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/api/secret/delete", r.repoAddress), bodyIOReader)
	if err != nil {
		return fmt.Errorf("err on DeleteSecret :%w", err)
	}
	request.Header.Add("Content-Type", "application/json")
	// request.Header.Add("X-Real-IP", "application/json")

	bearer := "Bearer " + user.Token
	request.Header.Add("Authorization", bearer)
	response, err := r.client.Do(request)
	if err != nil {
		return fmt.Errorf("err on DeleteSecret :%w", err)
	}
	// печатаем код ответа
	defer response.Body.Close()
	// читаем поток из тела ответа
	if response.StatusCode > 200 {
		if response.StatusCode == http.StatusNoContent {
			return fmt.Errorf("err on DeleteSecret :%w", secret.ErrSecretNotFound)
		}
		return fmt.Errorf("err on DeleteSecret: %w", fmt.Errorf(response.Status))
	}
	return nil
}

type getResponseSingle struct {
	Secret *Secret `json:"secret"`
}

func (r *SecretRequests) GetSecret(ctx context.Context, user *models.User, key string) (*models.Secret, error) {
	getInput := new(Input)
	getInput.Key = key
	bodyBytes, err := json.Marshal(getInput)
	if err != nil {
		return nil, fmt.Errorf("err on GetSecret:  %w", err)
	}
	bodyIOReader := bytes.NewReader(bodyBytes)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/api/secret/get", r.repoAddress), bodyIOReader)
	if err != nil {
		return nil, fmt.Errorf("err on GetSecret:  %w", err)
	}
	request.Header.Add("Content-Type", "application/json")
	// request.Header.Add("X-Real-IP", "application/json")

	bearer := "Bearer " + user.Token
	request.Header.Add("Authorization", bearer)
	response, err := r.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("err on GetSecret:  %w", err)
	}
	// печатаем код ответа
	defer response.Body.Close()
	// читаем поток из тела ответа
	if response.StatusCode > 200 {
		if response.StatusCode == http.StatusNoContent {
			return nil, fmt.Errorf("err on GetSecret: %w", secret.ErrSecretNotFound)
		}
		return nil, fmt.Errorf("err on GetSecret:  %w", fmt.Errorf(response.Status))
	}
	secretBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("err on GetSecret:  %w", err)
	}
	respSecret := new(getResponseSingle)
	err = json.Unmarshal(secretBody, respSecret)
	if err != nil {
		return nil, fmt.Errorf("err on GetSecret:  %w", err)
	}
	return &models.Secret{
		ID:     respSecret.Secret.ID,
		UserID: respSecret.Secret.UserID,
		Value:  respSecret.Secret.Value,
		Key:    key,
		Type:   respSecret.Secret.Type,
	}, nil
}

func (r *SecretRequests) GetSecrets(ctx context.Context, user *models.User) ([]*models.Secret, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/api/secret/list", r.repoAddress), nil)
	if err != nil {
		return nil, fmt.Errorf("err on GetSecrets :%w", err)
	}
	// request.Header.Add("X-Real-IP", "application/json")

	bearer := "Bearer " + user.Token
	request.Header.Add("Authorization", bearer)
	response, err := r.client.Do(request)
	if err != nil {

		return nil, fmt.Errorf("err on GetSecrets :%w", err)
	}
	// печатаем код ответа
	defer response.Body.Close()
	// читаем поток из тела ответа
	secretsBody, err := io.ReadAll(response.Body)
	if err != nil {

		return nil, fmt.Errorf("err on GetSecrets :%w", err)
	}
	if response.StatusCode > 200 {
		if response.StatusCode == http.StatusNoContent {
			return nil, fmt.Errorf("err on GetSecrets :%w", secret.ErrUserHaveNotSecret)
		}
		return nil, fmt.Errorf("err on GetSecrests :%w", fmt.Errorf(response.Status))
	}
	respSecrets := new(Secrets)
	err = json.Unmarshal(secretsBody, respSecrets)
	if err != nil {

		return nil, fmt.Errorf("err on GetSecrets :%w", err)
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
