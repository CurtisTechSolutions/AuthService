package internal

import (
	"log/slog"
	"os"
)

var logger *TLogger

type TLogger struct {
	Logger *slog.Logger
}

// InitializeLogger sets up the logger for the application.
// It configures the logger to output JSON formatted logs to stdout
func InitializeLogger(options *slog.HandlerOptions) {
	// Set up logger
	logger = &TLogger{
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, options)),
	}
}

func GetLogger() *slog.Logger {
	if logger == nil {
		InitializeLogger(&slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		})
	}
	return logger.Logger
}
