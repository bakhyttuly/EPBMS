package logger

import (
	"log/slog"
	"os"
)

// New creates a structured JSON logger for production or a text logger for development.
func New() *slog.Logger {
	env := os.Getenv("APP_ENV")
	if env == "development" {
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}
