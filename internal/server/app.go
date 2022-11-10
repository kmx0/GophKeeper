package server

import (
	"net/http"

	"github.com/kmx0/GophKeeper/internal/auth"
)

type App struct{
	httpServer *http.Server
	authUC auth.UseCase
}
func NewApp()*App{
	db := initDB()
	userRepo := auth
}