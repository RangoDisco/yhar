package config

import (
	"log/slog"
	"os"
)

func SetupLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo}))
}
