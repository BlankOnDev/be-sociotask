package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/harundarat/be-socialtask/internal/api"
	"github.com/harundarat/be-socialtask/internal/store"
	"github.com/harundarat/be-socialtask/migrations"
)

type Application struct {
	Logger      *log.Logger
	TaskHandler *api.TaskHandler
	UserHandler *api.UserHandler
	DB          *sql.DB
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

	// stores
	taskStore := store.NewPostgresTaskStore(pgDB)
	userStore := store.NewPostgresUserStore(pgDB)

	// handlers
	taskHandler := api.NewTaskHandler(taskStore, logger)
	userHandler := api.NewUserHandler(userStore, logger)

	app := &Application{
		Logger:      logger,
		TaskHandler: taskHandler,
		UserHandler: userHandler,
		DB:          pgDB,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}
