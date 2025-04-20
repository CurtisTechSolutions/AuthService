package main

import (
	"log/slog"
	"os"

	"github.com/CTS/AuthService/db"
	"github.com/CTS/AuthService/internal"
	"github.com/CTS/AuthService/server"
)

func Development() {
	// Initialize the database connection
	db.Connect(true)

	slog.Info("Starting server in development mode")
	// Initialize the logger
	internal.InitializeLogger(&slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})

}

func Production() {
	// Initialize the database connection
	db.Connect(false)

	slog.Info("Starting server in production mode")
	// Initialize the logger
	internal.InitializeLogger(&slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})

}

func main() {

	// Check if the environment is set to development or production
	env := os.Getenv("ENV")
	if env == "DEVELOPMENT" {
		Development()
	} else if env == "PRODUCTION" {
		Production()
	} else {
		slog.Error("Invalid environment variable", "ENV", env)
		slog.Warn("Defaulting to production mode")
		Production()
	}

	// Start the server
	err := server.Start()
	if err != nil {
		slog.Error("Server stopped", "error", err.Error())
		os.Exit(1)
	}
}
