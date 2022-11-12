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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	authhttp "github.com/kmx0/GophKeeper/internal/auth/delivery/http"
	secrethttp "github.com/kmx0/GophKeeper/internal/secret/delivery/http"
	scusecase "github.com/kmx0/GophKeeper/internal/secret/usecase"
	authmongo "github.com/zhashkevych/go-clean-architecture/auth/repository/mongo"
	authusecase "github.com/zhashkevych/go-clean-architecture/auth/usecase"
	bmmongo "github.com/zhashkevych/go-clean-architecture/bookmark/repository/mongo"
)

type App struct {
	httpServer *http.Server

	secretUC secret.UseCase
	authUC   auth.UseCase
}

func NewApp() *App {
	db := initDB()

	userRepo := authmongo.NewUserRepository(db, viper.GetString("mongo.user_collection"))
	secretkRepo := bmmongo.NewBookmarkRepository(db, viper.GetString("mongo.secret_collection"))

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

	secrethttp.RegisterHTTPEndpoints(api, a.secretUC)

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

func initDB() *mongo.Database {
	client, err := mongo.NewClient(options.Client().ApplyURI(viper.GetString("mongo.uri")))
	if err != nil {
		log.Fatalf("Error occured while establishing connection to mongoDB")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client.Database(viper.GetString("mongo.name"))
}
func initDBPG() {
	databaseUrl := "postgres://postgres:mypassword@localhost:5432/postgres"

	// this returns connection pool
	dbPool, err := pgxpool.Connect(context.Background(), databaseUrl)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	// to close DB pool
	defer dbPool.Close()
}
