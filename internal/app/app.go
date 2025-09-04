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

	// handlers
	taskHandler := api.NewTaskHandler(taskStore, logger)

	app := &Application{
		Logger:      logger,
		TaskHandler: taskHandler,
		DB:          pgDB,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}
