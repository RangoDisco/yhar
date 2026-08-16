package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/getsentry/sentry-go"
	sentryslog "github.com/getsentry/sentry-go/slog"
	"github.com/gin-gonic/gin"
)

func SetupSentry() {
	sentryDSN := os.Getenv("SENTRY_DSN")
	if sentryDSN != "" {
		// To initialize Sentry's handler, you need to initialize Sentry itself beforehand
		if err := sentry.Init(sentry.ClientOptions{
			Dsn: sentryDSN,
			DataCollection: &sentry.DataCollection{
				UserInfo: sentry.Set(true),
			},
		}); err != nil {
			fmt.Printf("Sentry initialization failed: %v\n", err)
		}
	}
}

func SetupLogger(ctx context.Context) {
	var minLevel slog.Level
	var logger *slog.Logger

	switch gin.Mode() {
	case gin.DebugMode:
		minLevel = slog.LevelDebug
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: minLevel}))
	default:
		handler := sentryslog.Option{
			LogLevel:  []slog.Level{slog.LevelInfo, slog.LevelWarn, slog.LevelWarn, slog.LevelError, sentryslog.LevelFatal},
			AddSource: true,
		}.NewSentryHandler(ctx)
		logger = slog.New(handler)
	}
	slog.SetDefault(logger)
}
