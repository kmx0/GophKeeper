package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kmx0/GophKeeper/config"
	authcli "github.com/kmx0/GophKeeper/internal/auth/delivery/cli"
	authremotestorage "github.com/kmx0/GophKeeper/internal/auth/repository/remote"
	authremoteusecase "github.com/kmx0/GophKeeper/internal/auth/usecase/remote"
	secretcli "github.com/kmx0/GophKeeper/internal/secret/delivery/cli"
	secretremotestorage "github.com/kmx0/GophKeeper/internal/secret/repository/remote"
	secretremoteusecase "github.com/kmx0/GophKeeper/internal/secret/usecase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {

	if err := config.Init(); err != nil {
		log.Fatalf("%s", err.Error())
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
	userRepo := authremotestorage.NewUserRepository(&http.Client{Transport: tr}, fmt.Sprintf("https://%s:%s", viper.GetString("address"), viper.GetString("port")))
	secretRepo := secretremotestorage.NewSecretRepository(&http.Client{Transport: tr}, fmt.Sprintf("https://%s:%s", viper.GetString("address"), viper.GetString("port")))

	authcli.RegisterAuthCmdEndpoints(rootCmd, authremoteusecase.NewAuthUseCase(
		userRepo,
		viper.GetString("auth.hash_salt"),
		[]byte(viper.GetString("auth.signing_key")),
		viper.GetString("auth.token_file"),
	), os.Stdout)

	authMidleWare := authcli.NewAuthMiddleware(authremoteusecase.NewAuthUseCase(userRepo,
		viper.GetString("auth.hash_salt"),
		[]byte(viper.GetString("auth.signing_key")),
		viper.GetString("auth.token_file"),
	),
		viper.GetString("auth.token_file"))
	secretcli.RegisterSecretCmdEndpoints(rootCmd, secretremoteusecase.NewSecretUseCase(secretRepo), authMidleWare)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

//TODO
//base64 <->ок
//сохранение, если указан путь к файлу ок
//app -> client
