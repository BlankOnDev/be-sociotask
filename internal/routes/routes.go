package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/harundarat/be-socialtask/internal/app"
	"github.com/harundarat/be-socialtask/internal/pages"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(pages.IndexPage))
	})

	r.Get("/success", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(pages.SuccessPage))
	})

	r.Get("/failed", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(pages.FailedPage))
	})

	// r.Group(func(r chi.Router) {
	// 	r.Use(app.UserMiddleware.Authenticate)
	// })

	r.Get("/health", app.HealthCheck)
	r.Get("/tasks/{id}", app.TaskHandler.HandleGetTaskByID)
	r.Get("/users/{id}/tasks", app.UserHandler.HandleGetUserTasks)
	r.Get("/login/twitter", app.AuthHandler.HandleTwitterLogin)
	r.Get("/login/twitter/callback", app.AuthHandler.HandleTwitterCallback)
	r.Post("/tasks", app.TaskHandler.HandleCreateTask)
	r.Post("/register", app.UserHandler.HandleCreateUser)
	r.Post("/login", app.UserHandler.HandleLoginUser)
	r.Get("/login/google", app.AuthHandler.LoginAuthenticationGooogle)
	r.Get("/login/google/callback", app.AuthHandler.CallbackAuthenticationGooogle)
	r.Post("/login/google/android", app.AuthHandler.HandleGoogleLoginAndroid)

	r.Group(func(r chi.Router) {
		r.Use(app.UserMiddleware.Authenticate)
		r.Use(func(next http.Handler) http.Handler {
			return app.UserMiddleware.RequireUser(next.ServeHTTP)
		})

		r.Get("/tasks", app.TaskHandler.HandleGetAllTask)
		r.Post("/tasks", app.TaskHandler.HandleCreateTask)
		r.Put("/tasks/{id}", app.TaskHandler.HandleEditTask)
		r.Delete("/tasks/{id}", app.TaskHandler.HandleDeleteTask)
	})

	return r
}
