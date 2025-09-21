package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/harundarat/be-socialtask/internal/api"
	"github.com/harundarat/be-socialtask/internal/middleware"
	"github.com/harundarat/be-socialtask/internal/store"
	"github.com/harundarat/be-socialtask/migrations"
	"golang.org/x/oauth2"
)

type Application struct {
	Logger         *log.Logger
	TaskHandler    *api.TaskHandler
	UserHandler    *api.UserHandler
	AuthHandler    *api.AuthHandler
	UserMiddleware *middleware.UserMiddleware
	DB             *sql.DB
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	oauthConf := &oauth2.Config{
		ClientID:     store.GetEnv("TWITTER_CLIENT_ID"),
		ClientSecret: store.GetEnv("TWITTER_CLIENT_SECRET"),
		RedirectURL:  store.GetEnv("TWITTER_REDIRECT_URL"),
		Scopes:       []string{"tweet.read", "users.read", "offline.access"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://twitter.com/i/oauth2/authorize",
			TokenURL: "https://api.twitter.com/2/oauth2/token",
		},
	}

	// stores
	taskStore := store.NewPostgresTaskStore(pgDB)
	userStore := store.NewPostgresUserStore(pgDB)

	// handlers
	taskHandler := api.NewTaskHandler(taskStore, logger)
	userHandler := api.NewUserHandler(userStore, logger)
	authHandler := api.NewAuthHandler(logger, userStore, oauthConf)

	// middleware
	userMiddleware := middleware.NewUserMiddleware(userStore, "thisissecret")

	app := &Application{
		Logger:         logger,
		TaskHandler:    taskHandler,
		UserHandler:    userHandler,
		AuthHandler:    authHandler,
		UserMiddleware: userMiddleware,
		DB:             pgDB,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}
