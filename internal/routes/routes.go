package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/harundarat/be-socialtask/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(app.UserMiddleware.Authenticate)
	})

	r.Get("/health", app.HealthCheck)
	r.Get("/tasks/{id}", app.TaskHandler.HandleGetTaskByID)
	r.Post("/tasks", app.TaskHandler.HandleCreateTask)
	r.Post("/users", app.UserHandler.HandleCreateUser)
	r.Post("/login", app.UserHandler.HandleLoginUser)

	return r
}
