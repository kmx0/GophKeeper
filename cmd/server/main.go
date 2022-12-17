package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/kmx0/GophKeeper/config"
	"github.com/kmx0/GophKeeper/internal/server"
	"github.com/spf13/viper"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("error on running server: %v", err.Error())
	}

	fulladdr := viper.GetString("gokeeper_server_addr")
	temp := strings.Split(fulladdr, "//")
	if len(temp) < 2 {
		log.Fatalf("error on running server: %v", fmt.Errorf("incorrect addr in config file, need proto://addr:port"))
	}
	addr := temp[1]
	app := server.NewApp(context.Background(), viper.GetString("dsn"))
	if err := app.Run(addr); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
