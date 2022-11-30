package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kmx0/GophKeeper/config"
	authcli "github.com/kmx0/GophKeeper/internal/auth/delivery/cli"
	authremotestorage "github.com/kmx0/GophKeeper/internal/auth/repository/remote"
	authrequests "github.com/kmx0/GophKeeper/internal/auth/repository/remote/requests"
	authremoteusecase "github.com/kmx0/GophKeeper/internal/auth/usecase/remote"
	secretcli "github.com/kmx0/GophKeeper/internal/secret/delivery/cli"
	secretremotestorage "github.com/kmx0/GophKeeper/internal/secret/repository/remote"
	secretrequests "github.com/kmx0/GophKeeper/internal/secret/repository/remote/requests"
	secretremoteusecase "github.com/kmx0/GophKeeper/internal/secret/usecase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {

	if err := config.Init(); err != nil {
		log.Fatalf("error on client running: %v", err.Error())
	}
	rootCmd := &cobra.Command{
		Use:     "gophkeeper",
		Short:   "Keep your data securly 💻",
		Version: "0.1",
	}
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	}

	userRequests := authrequests.NewUserRequests(&http.Client{Transport: tr}, viper.GetString("gokeeper_server_addr"))
	userRepo := authremotestorage.NewUserRepository(userRequests)
	secretResquests := secretrequests.NewSecretRequest(&http.Client{Transport: tr}, viper.GetString("gokeeper_server_addr"))
	secretRepo := secretremotestorage.NewSecretRepository(secretResquests)

	authcli.RegisterAuthCmdEndpoints(rootCmd, authremoteusecase.NewAuthUseCase(
		userRepo,
		viper.GetString("auth.hash_salt"),
		[]byte(viper.GetString("auth.signing_key")),
		viper.GetString("auth.token_file"),
	), os.Stdout)

	authStatus := authcli.NewAuthStatus(authremoteusecase.NewAuthUseCase(userRepo,
		viper.GetString("auth.hash_salt"),
		[]byte(viper.GetString("auth.signing_key")),
		viper.GetString("auth.token_file"),
	),
		viper.GetString("auth.token_file"))
	secretcli.RegisterSecretCmdEndpoints(rootCmd, secretremoteusecase.NewSecretUseCase(secretRepo), authStatus, os.Stdout)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
