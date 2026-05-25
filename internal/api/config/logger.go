package config

import (
	"log/slog"
	"os"
)

func SetupLogger() {
	var minLevel slog.Level

	switch os.Getenv("APP_ENV") {
	case "debug":
		minLevel = slog.LevelDebug
	default:
		minLevel = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: minLevel}))
	slog.SetDefault(logger)
}
