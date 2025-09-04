package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/harundarat/be-socialtask/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)
	r.Get("/task/{id}", app.TaskHandler.HandleGetTaskByID)
	r.Post("/task", app.TaskHandler.HandleCreateTask)

	return r
}
