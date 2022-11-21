package main

import (
	"context"
	"log"

	"github.com/kmx0/GophKeeper/config"
	"github.com/kmx0/GophKeeper/internal/server"
	"github.com/spf13/viper"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	app := server.NewApp(context.Background(), "")

	if err := app.Run(viper.GetString("port")); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
//pg migrate
//tls
//pg