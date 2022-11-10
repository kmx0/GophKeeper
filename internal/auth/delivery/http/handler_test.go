package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kmx0/GophKeeper/internal/auth/usecase"
)

func TestSignUp(t *testing.T) {
	r := gin.Default()
	uc := new(usecase.AuthUseCaseMock)
	RegisterHTTPEndpoints(r, uc)

	signUpBody := &signInput{
		Login:    "testuser",
		Password: "testpass",
	}
	body, err := json.Marshal(signUpBody)
	assert.NoError(t, err)
	uc.On("SignUp", signUpBody.Login, signUpBody.Password).Return(nil)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/sign-up", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

}
