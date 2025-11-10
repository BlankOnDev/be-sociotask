package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/harundarat/be-socialtask/internal/app"
	"github.com/harundarat/be-socialtask/internal/pages"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

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
	r.Get("/tasks", app.TaskHandler.HandleGetAllTask)
	r.Get("/tasks/{id}", app.TaskHandler.HandleGetTaskByID)
	r.Get("/users/{id}/tasks", app.UserHandler.HandleGetUserTasks)
	r.Get("/login/twitter", app.AuthHandler.HandleTwitterLogin)
	r.Get("/login/twitter/callback", app.AuthHandler.HandleTwitterCallback)
	r.Post("/register", app.UserHandler.HandleCreateUser)
	r.Post("/login", app.UserHandler.HandleLoginUser)
	r.Get("/login/google", app.AuthHandler.LoginAuthenticationGooogle)
	r.Get("/login/google/callback", app.AuthHandler.CallbackAuthenticationGooogle)
	r.Post("/login/google/android", app.AuthHandler.HandleGoogleLoginAndroid)

	// task actions
	r.Get("/actions", app.ActionHandler.HandleGetAllAction)
	r.Get("/actions/{id}", app.ActionHandler.HandleGetActionByID)

	// task rewards
	r.Get("/reward", app.RewardHandler.HandleGetAllReward)
	r.Get("/reward/{id}", app.RewardHandler.HandleGetRewardByID)

	r.Group(func(r chi.Router) {
		r.Use(app.UserMiddleware.Authenticate)
		r.Use(func(next http.Handler) http.Handler {
			return app.UserMiddleware.RequireUser(next.ServeHTTP)
		})

		// user
		r.Get("/users/current", app.UserMiddleware.RequireUser(app.UserHandler.HandleGetCurrentUser))

		// task
		r.Post("/tasks", app.TaskHandler.HandleCreateTask)
		r.Put("/tasks/{id}", app.TaskHandler.HandleEditTask)
		r.Delete("/tasks/{id}", app.TaskHandler.HandleDeleteTask)

		// // action
		// r.Post("/actions", app.ActionHandler.HandleCreateAction)
		// r.Put("/actions/{id}", app.ActionHandler.HandleEditAction)
		// r.Delete("/actions/{id}", app.ActionHandler.HandleDeleteAction)

		// // reward
		// r.Post("/reward", app.RewardHandler.HandleCreateReward)
		// r.Put("/reward/{id}", app.RewardHandler.HandleEditReward)
		// r.Delete("/reward/{id}", app.RewardHandler.HandleDeleteReward)

		// rewards
		r.Post("/rewards", app.RewardsHandler.HandleCreateReward)
	})

	return r
}
