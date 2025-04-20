package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func Start() error {
	// Initialize the server
	r := chi.NewRouter()

	// Set up routes
	// Define middleware
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(middleware.Logger) // <--<< Logger should come before Recoverer
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Heartbeat("/healthcheck"))
	r.Use(middleware.Recoverer)

	// Define routes
	r.Mount("/", BaseRoutes())

	// Start the server
	err := http.ListenAndServe(":9090", r)
	return err
}
