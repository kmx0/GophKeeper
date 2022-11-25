package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kmx0/GophKeeper/internal/auth/usecase"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestSignIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	uc := usecase.NewMockUseCase(ctrl)

	// гарантируем, что заглушка
	// при вызове с аргументом "Key" вернёт "Value"
	// login = ""
	rootCmd := &cobra.Command{
		Use:     "gophkeeper",
		Short:   "Keep your data securly 💻",
		Version: "0.1",
	}
	b := bytes.NewBufferString("")
	// rootCmd.SetOut(b)
	RegisterAuthCmdEndpoints(rootCmd, uc, b)
	tests := []struct {
		name    string
		command string
		token   string
		err     error
		want    error
	}{

		{
			name:    "internal server error sign-up",
			command: "sign-up --login=log --password=pas",
			err:     nil,
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc.EXPECT().SignUp(ctx, login, password).Return(tt.err).AnyTimes()

			rootCmd.SetArgs(strings.Split(tt.command, " "))
			err := rootCmd.Execute()
			assert.Equal(t, tt.want, err)
		})
	}

}

// func TestSignUp(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	ctx := context.Background()
// 	uc := usecase.NewMockUseCase(ctrl)

// 	// гарантируем, что заглушка
// 	// при вызове с аргументом "Key" вернёт "Value"
// 	r := gin.Default()
// 	RegisterHTTPEndpoints(r, uc)

// 	//->
// 	// тестируем функцию
// 	tests := []struct {
// 		name        string
// 		corruptBody bool
// 		creds       *signInput
// 		err         error
// 		want        int
// 	}{
// 		{
// 			name: "correct sign-up",
// 			creds: &signInput{
// 				Login:    "testuser",
// 				Password: "testpass",
// 			},
// 			err:  nil,
// 			want: 200,
// 		},
// 		{
// 			name: "conflict sign-up",
// 			creds: &signInput{
// 				Login:    "testuser",
// 				Password: "testpass",
// 			},
// 			err:  auth.ErrLoginBusy,
// 			want: 409,
// 		},
// 		{
// 			name: "corrupt body",
// 			creds: &signInput{
// 				Login:    "1w2",
// 				Password: "212",
// 			},
// 			err:         fmt.Errorf("err"),
// 			want:        400,
// 			corruptBody: true,
// 		},
// 		{
// 			name: "internal server error sign-up",
// 			creds: &signInput{
// 				Login:    "1",
// 				Password: "1",
// 			},
// 			err:  fmt.Errorf("internal server error"),
// 			want: 500,
// 		},
// 	}
// 	var calls []*gomock.Call
// 	for _, tt := range tests {

// 		var call *gomock.Call
// 		if tt.corruptBody {
// 			call = uc.EXPECT().SignUp(ctx, tt.creds.Login, tt.creds.Password).Return(tt.err).AnyTimes()
// 		} else {
// 			call = uc.EXPECT().SignUp(ctx, tt.creds.Login, tt.creds.Password).Return(tt.err)
// 		}
// 		calls = append(calls, call)
// 	}
// 	gomock.InOrder(calls...)
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var body []byte
// 			var err error
// 			if tt.corruptBody {
// 				// body, err = json.Marshal(tt.corruptCreds)
// 				body, err = json.Marshal(tt.creds)
// 				body = body[:len(body)-2]
// 			} else {
// 				body, err = json.Marshal(tt.creds)
// 			}
// 			assert.NoError(t, err)
// 			w := httptest.NewRecorder()
// 			req, _ := http.NewRequest("POST", "/auth/sign-up", bytes.NewBuffer(body))
// 			r.ServeHTTP(w, req)
// 			assert.Equal(t, tt.want, w.Code)
// 		})
// 	}

// 	//->
// 	// w = httptest.NewRecorder()
// 	// req, _ = http.NewRequest("POST", "/auth/sign-up", bytes.NewBuffer(body))
// 	// r.ServeHTTP(w, req)
// 	// assert.Equal(t, 409, w.Code)

// }
