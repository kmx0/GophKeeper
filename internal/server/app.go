package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kmx0/GophKeeper/internal/auth"
	"github.com/kmx0/GophKeeper/internal/secret"
	"github.com/spf13/viper"

	authhttp "github.com/kmx0/GophKeeper/internal/auth/delivery/http"
	authlocalstorage "github.com/kmx0/GophKeeper/internal/auth/repository/localstorage"
	authpostgres "github.com/kmx0/GophKeeper/internal/auth/repository/postgres"
	authusecase "github.com/kmx0/GophKeeper/internal/auth/usecase"
	schttp "github.com/kmx0/GophKeeper/internal/secret/delivery/http"
	sclocalstorage "github.com/kmx0/GophKeeper/internal/secret/repository/localstorage"
	scpostgres "github.com/kmx0/GophKeeper/internal/secret/repository/postgres"
	scusecase "github.com/kmx0/GophKeeper/internal/secret/usecase"
)

type App struct {
	httpServer *http.Server

	secretUC secret.UseCase
	authUC   auth.UseCase
}

func NewApp(ctx context.Context, dsn string) *App {
	if dsn != "" {
		pool, err := initDBPG(ctx, dsn)
		if err != nil {
			return nil
		}
		userRepo := authpostgres.NewUserRepository(pool)
		secretkRepo := scpostgres.NewSecretRepository(pool)
		return &App{
			secretUC: scusecase.NewSecretUseCase(secretkRepo),
			authUC: authusecase.NewAuthUseCase(
				userRepo,
				viper.GetString("auth.hash_salt"),
				[]byte(viper.GetString("auth.signing_key")),
				viper.GetDuration("auth.token_ttl"),
			),
		}
	} else {
		userRepo := authlocalstorage.NewUserLocalStorage()
		secretkRepo := sclocalstorage.NewSecretLocalStorage()

		return &App{
			secretUC: scusecase.NewSecretUseCase(secretkRepo),
			authUC: authusecase.NewAuthUseCase(
				userRepo,
				viper.GetString("auth.hash_salt"),
				[]byte(viper.GetString("auth.signing_key")),
				viper.GetDuration("auth.token_ttl"),
			),
		}
	}

}

func (a *App) Run(port string) error {
	// Init gin handler
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	// Set up http handlers
	// SignUp/SignIn endpoints
	authhttp.RegisterHTTPEndpoints(router, a.authUC)

	// API endpoints
	authMiddleware := authhttp.NewAuthMiddleware(a.authUC)
	api := router.Group("/api", authMiddleware)

	schttp.RegisterHTTPEndpoints(api, a.secretUC)

	// HTTP Server
	a.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}

func initDBPG(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	conf, err := pgxpool.ParseConfig(dsn) // Using environment variables instead of a connection string.
	if err != nil {
		return nil, err
	}

	conf.LazyConnect = true

	pool, err := pgxpool.ConnectConfig(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("pgx connection error: %w", err)
	}

	return pool, nil
}
