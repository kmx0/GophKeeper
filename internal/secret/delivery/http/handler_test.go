package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kmx0/GophKeeper/internal/auth"
	"github.com/kmx0/GophKeeper/internal/models"
	"github.com/kmx0/GophKeeper/internal/secret/usecase"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	testUser := &models.User{
		Login:    "testuser",
		Password: "password",
	}
	r := gin.Default()

	group := r.Group("/api", func(c *gin.Context) {
		c.Set(auth.CtxUserKey, testUser)
	})
	uc := new(usecase.SecretUseCaseMock)

	RegisterHTTPEndpoints(group, uc)

	scs := make([]*models.Secret, 5)
	for i := 0; i < 5; i++ {
		scs[i] = &models.Secret{
			ID:     i,
			UserID: i * i,
			Key:    fmt.Sprintf("key%d", i),
			Value:  fmt.Sprintf("value%d", i),
		}
	}

	uc.On("GetSecrets", testUser).Return(scs, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/secret/list", nil)
	r.ServeHTTP(w, req)
	expectOut := &getResponse{Secrets: toSecrets(scs)}
	expectOutBody, err := json.Marshal(expectOut)
	assert.NoError(t, err)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, string(expectOutBody), w.Body.String())

}

func TestGet(t *testing.T) {
	testUser := &models.User{
		Login:    "testuser",
		Password: "password",
	}
	r := gin.Default()

	group := r.Group("/api", func(c *gin.Context) {
		c.Set(auth.CtxUserKey, testUser)
	})
	uc := new(usecase.SecretUseCaseMock)

	RegisterHTTPEndpoints(group, uc)

	sc := &models.Secret{
		ID:     1,
		UserID: 11,
		Key:    "key",
		Value:  "value",
	}
	inp := &getInput{
		Key: "key",
	}
	body, err := json.Marshal(inp)
	assert.NoError(t, err)

	uc.On("GetSecret", testUser, inp.Key).Return(sc, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/secret/get", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	expectOut := &getResponseSingle{Secret: toSecret(sc)}
	expectOutBody, err := json.Marshal(expectOut)
	assert.NoError(t, err)

	assert.Equal(t, 200, w.Code)

	assert.Equal(t, string(expectOutBody), w.Body.String())

}
