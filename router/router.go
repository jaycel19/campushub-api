package router

import (
	"net/http"

	"github.com/jaycel19/campushub-api/controllers"
	"github.com/jaycel19/campushub-api/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Routes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// AUTH ROUTES
	router.Post("/api/v1/users/login", controllers.UserLogin)
	router.Post("/api/v1/users", controllers.CreateUser)
	// TODO: Add route for refresh tokens

	router.Group(func(auth chi.Router) {
		auth.Use(middlewares.RequireAuth)

		auth.Post("/api/v1/users/login/renew_token", controllers.RenewToken)
		// USER ROUTES
		auth.Get("/api/v1/users", controllers.GetAllUser)
		auth.Get("/api/v1/users/{username}", controllers.GetUserById)

		// POST ROUTES
		auth.Get("/api/v1/posts", controllers.GetAllPosts)
		auth.Get("/api/v1/posts/{id}", controllers.GetPostById)
		auth.Post("/api/v1/posts", controllers.CreatePost)
		auth.Put("/api/v1/posts/{id}", controllers.UpdatePost)
		auth.Delete("/api/v1/posts/{id}", controllers.DeletePost)
	})

	return router
}
