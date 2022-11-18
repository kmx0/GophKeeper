package main

import (
	"log"
	"net/http"

	"github.com/kmx0/GophKeeper/internal/auth/delivery/cli"
	authremotestorage "github.com/kmx0/GophKeeper/internal/auth/repository/remote"
	authremoteusecase "github.com/kmx0/GophKeeper/internal/auth/usecase/remote"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "gophkeeper",
		Short:   "Keep your data securly 💻",
		Version: "0.1",
	}
	userRepo := authremotestorage.NewUserRepository(&http.Client{},"http://localhost:8000")

	cli.RegisterCmdEndpoints(rootCmd, authremoteusecase.NewAuthUseCase(
		userRepo,
		viper.GetString("auth.hash_salt"),
		[]byte(viper.GetString("auth.signing_key")),
	))

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
