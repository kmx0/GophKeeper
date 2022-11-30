package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kmx0/GophKeeper/internal/auth"
	"github.com/kmx0/GophKeeper/internal/auth/usecase"
	"github.com/stretchr/testify/assert"
)

func TestSignIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	uc := usecase.NewMockUseCase(ctrl)

	// гарантируем, что заглушка
	// при вызове с аргументом "Key" вернёт "Value"
	r := gin.Default()
	RegisterHTTPEndpoints(r, uc)
	tests := []struct {
		name        string
		corruptBody bool
		creds       *signInput
		token       string
		err         error
		want        int
	}{
		{
			name: "correct sign-in",
			creds: &signInput{
				Login:    "testuser",
				Password: "testpass",
			},
			err:   nil,
			token: "token",
			want:  200,
		},
		{
			name: "conflict sign-in",
			creds: &signInput{
				Login:    "incorrectuser",
				Password: "testpass",
			},
			err:  auth.ErrUserNotFound,
			want: 401,
		},
		{
			name: "corrupt body",
			creds: &signInput{
				Login:    "1w2",
				Password: "212",
			},
			err:         fmt.Errorf("err"),
			want:        400,
			corruptBody: true,
		},
		{
			name: "internal server error sign-in",
			creds: &signInput{
				Login:    "1",
				Password: "1",
			},
			err:  fmt.Errorf("internal server error"),
			want: 500,
		},
	}
	var calls []*gomock.Call
	for _, tt := range tests {

		var call *gomock.Call
		if tt.corruptBody {
			call = uc.EXPECT().SignIn(ctx, tt.creds.Login, tt.creds.Password).Return("", tt.err).AnyTimes()
		} else {
			call = uc.EXPECT().SignIn(ctx, tt.creds.Login, tt.creds.Password).Return(tt.token, tt.err)
		}
		calls = append(calls, call)
	}
	gomock.InOrder(calls...)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			var err error
			if tt.corruptBody {
				// body, err = json.Marshal(tt.corruptCreds)
				body, err = json.Marshal(tt.creds)
				body = body[:len(body)-2]
			} else {
				body, err = json.Marshal(tt.creds)
			}
			assert.NoError(t, err)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/auth/sign-in", bytes.NewBuffer(body))
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.want, w.Code)
		})
	}

}

func TestSignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	uc := usecase.NewMockUseCase(ctrl)

	// гарантируем, что заглушка
	// при вызове с аргументом "Key" вернёт "Value"
	r := gin.Default()
	RegisterHTTPEndpoints(r, uc)

	//->
	// тестируем функцию
	tests := []struct {
		name        string
		corruptBody bool
		creds       *signInput
		err         error
		want        int
	}{
		{
			name: "correct sign-up",
			creds: &signInput{
				Login:    "testuser",
				Password: "testpass",
			},
			err:  nil,
			want: 200,
		},
		{
			name: "conflict sign-up",
			creds: &signInput{
				Login:    "testuser",
				Password: "testpass",
			},
			err:  auth.ErrLoginBusy,
			want: 409,
		},
		{
			name: "corrupt body",
			creds: &signInput{
				Login:    "1w2",
				Password: "212",
			},
			err:         fmt.Errorf("err"),
			want:        400,
			corruptBody: true,
		},
		{
			name: "internal server error sign-up",
			creds: &signInput{
				Login:    "1",
				Password: "1",
			},
			err:  fmt.Errorf("internal server error"),
			want: 500,
		},
	}
	var calls []*gomock.Call
	for _, tt := range tests {

		var call *gomock.Call
		if tt.corruptBody {
			call = uc.EXPECT().SignUp(ctx, tt.creds.Login, tt.creds.Password).Return(tt.err).AnyTimes()
		} else {
			call = uc.EXPECT().SignUp(ctx, tt.creds.Login, tt.creds.Password).Return(tt.err)
		}
		calls = append(calls, call)
	}
	gomock.InOrder(calls...)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			var err error
			if tt.corruptBody {
				// body, err = json.Marshal(tt.corruptCreds)
				body, err = json.Marshal(tt.creds)
				body = body[:len(body)-2]
			} else {
				body, err = json.Marshal(tt.creds)
			}
			assert.NoError(t, err)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/auth/sign-up", bytes.NewBuffer(body))
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.want, w.Code)
		})
	}

	//->
	// w = httptest.NewRecorder()
	// req, _ = http.NewRequest("POST", "/auth/sign-up", bytes.NewBuffer(body))
	// r.ServeHTTP(w, req)
	// assert.Equal(t, 409, w.Code)

}
