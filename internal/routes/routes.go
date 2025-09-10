package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/harundarat/be-socialtask/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)
	r.Get("/tasks/{id}", app.TaskHandler.HandleGetTaskByID)
	r.Post("/tasks", app.TaskHandler.HandleCreateTask)
	r.Post("/users", app.UserHandler.HandleCreateUser)

	return r
}
