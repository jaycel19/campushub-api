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
	router.Get("/api/v1/users/{username}", controllers.GetUserById)
	// TODO: Add route for refresh tokens
	router.Post("/api/v1/users/login/renew_token", controllers.RenewToken)

	router.Group(func(auth chi.Router) {
		auth.Use(middlewares.RequireAuth)

		// USER ROUTES
		auth.Get("/api/v1/users", controllers.GetAllUser)

		// PROFILE ROUTEs
		auth.Get("/api/v1/profiles", controllers.GetAllProfiles)
		auth.Get("/api/v1/profiles/{username}", controllers.GetProfileByUser)
		auth.Put("/api/v1/profiles/{username}", controllers.UpdateProfile)
		auth.Post("/api/v1/profiles", controllers.CreateProfile)
		auth.Post("/api/v1/profiles/{username}/background", controllers.ProfileChangeBackground)

		// POST ROUTES
		auth.Get("/api/v1/posts", controllers.GetAllPosts)
		auth.Get("/api/v1/posts/{id}", controllers.GetPostById)
		auth.Post("/api/v1/posts/{id}/like", controllers.PostLike)
		auth.Post("/api/v1/posts", controllers.CreatePost)
		auth.Put("/api/v1/posts/{id}", controllers.UpdatePost)
		auth.Delete("/api/v1/posts/{id}", controllers.DeletePost)

		// COMMENTS ROUTES
		auth.Get("/api/v1/comments", controllers.GetAllComments)
		auth.Post("/api/v1/comments", controllers.CreateComment)
		auth.Get("/api/v1/comments/{pid}", controllers.GetCommentsByPostID)
		auth.Put("/api/v1/comments", controllers.UpdateComment)
		auth.Delete("/api/v1/comments/{id}", controllers.DeleteComment)
	})

	return router
}
