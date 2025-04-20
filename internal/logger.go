package internal

import (
	"log/slog"
	"os"
)

// InitializeLogger sets up the logger for the application.
// It configures the logger to output JSON formatted logs to stdout
func InitializeLogger(options *slog.HandlerOptions) {
	// Set up logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, options))
	slog.SetDefault(logger)
}
