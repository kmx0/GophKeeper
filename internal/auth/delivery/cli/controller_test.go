package cli

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kmx0/GophKeeper/internal/auth"
	"github.com/kmx0/GophKeeper/internal/auth/usecase"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestSignUp(t *testing.T) {
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
	type creds struct {
		login    string
		password string
	}
	tests := []struct {
		name     string
		command  string
		creds    creds
		token    string
		err      error
		want     string
		dontCall bool
	}{

		{
			name:     "empty login",
			command:  "sign-up --login= --password=pas",
			creds:    creds{login: "", password: "pas"},
			err:      nil,
			want:     "login is empty!",
			dontCall: true,
		},
		{
			name:     "empty password",
			command:  "sign-up --login=mylogin --password=",
			creds:    creds{login: "mylogin", password: ""},
			err:      nil,
			want:     "password is empty!",
			dontCall: true,
		},
		{
			name:    "internal server",
			command: "sign-up --login=mylogin --password=password",
			creds:   creds{login: "mylogin", password: "password"},
			err:     fmt.Errorf("internal server error"),
			want:    "internal server error: internal server error",
		},
		{
			name:    "success",
			command: "sign-up --login=mylogin --password=password",
			creds:   creds{login: "mylogin", password: "password"},
			err:     nil,
			want:    "Successfully registered",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dontCall {
				uc.EXPECT().SignUp(ctx, tt.creds.login, tt.creds.password).Return(tt.err).AnyTimes()
			} else {
				uc.EXPECT().SignUp(ctx, tt.creds.login, tt.creds.password).Return(tt.err)

			}
			rootCmd.SetArgs(strings.Split(tt.command, " "))
			rootCmd.Execute()
			out, err := ioutil.ReadAll(b)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.want, string(out))
		})

	}

}

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
	type creds struct {
		login    string
		password string
	}
	tests := []struct {
		name     string
		command  string
		creds    creds
		token    string
		err      error
		want     string
		dontCall bool
	}{

		{
			name:     "empty login",
			command:  "sign-in --login= --password=pas",
			creds:    creds{login: "", password: "pas"},
			err:      nil,
			want:     "login is empty!",
			dontCall: true,
		},
		{
			name:     "empty password",
			command:  "sign-in --login=mylogin --password=",
			creds:    creds{login: "mylogin", password: ""},
			err:      nil,
			want:     "password is empty!",
			dontCall: true,
		},
		{
			name:    "user not found",
			command: "sign-in --login=incorrect --password=incorrect",
			creds:   creds{login: "incorrect", password: "incorrect"},
			err:     auth.ErrUserNotFound,
			want:    "not such login or password",
		},
		{
			name:    "internal server",
			command: "sign-in --login=mylogin --password=password",
			creds:   creds{login: "mylogin", password: "password"},
			err:     fmt.Errorf("internal server error"),
			want:    "internal server error: internal server error",
		},
		{
			name:    "success",
			command: "sign-in --login=mylogin --password=password",
			creds:   creds{login: "mylogin", password: "password"},
			err:     nil,
			want:    "Successfully logged in",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dontCall {
				uc.EXPECT().SignIn(ctx, tt.creds.login, tt.creds.password).Return(tt.token, tt.err).AnyTimes()
			} else {
				uc.EXPECT().SignIn(ctx, tt.creds.login, tt.creds.password).Return(tt.token, tt.err)

			}
			rootCmd.SetArgs(strings.Split(tt.command, " "))
			rootCmd.Execute()
			out, err := ioutil.ReadAll(b)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.want, string(out))
		})

	}

}
