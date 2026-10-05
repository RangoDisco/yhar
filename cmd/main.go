package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/api/config"
	"github.com/rangodisco/yhar/api/middlewares"
)

func init() {
	if os.Getenv("JWT_SECRET") == "" {
		log.Fatalf("JWT_SECRET environment variable not set")
	}

	if os.Getenv("BASE_URL") == "" {
		log.Fatalf("BASE_URL environment variable not set")
	}

	if os.Getenv("GIN_MODE") == "" {
		log.Fatalf("GIN_MODE environment variable not set")
	}
}

func main() {
	yDb, err := config.SetupDatabase()
	if err != nil {
		log.Fatalf("failed to init database: %v", err)
	}

	switch os.Getenv("GIN_MODE") {
	case gin.DebugMode:
		gin.SetMode(gin.DebugMode)
	case gin.TestMode:
		gin.SetMode(gin.TestMode)
	case gin.ReleaseMode:
		gin.SetMode(gin.ReleaseMode)
		config.SetupSentry()
		// Flush buffered events before the program terminates.
		// Set the timeout to the maximum duration the program can afford to wait.
		defer sentry.Flush(10 * time.Second)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	_, serverServices, handlers, pollers, importers := config.AutoWire(yDb)

	// Start pollers (subsonic only for now)
	if os.Getenv("GIN_MODE") == gin.ReleaseMode && pollers.Subsonic != nil {
		go pollers.Subsonic.Start(ctx)
		log.Println("Started Subsonic poller")
	}

	// Import historic data
	if importers.Maloja != nil {
		go func() {
			err := importers.Maloja.Import(ctx)
			if err != nil {
				log.Printf("Failed to import Maloja's data :%v", err)
			}
		}()
	}

	config.SetupLogger(ctx)
	// Start router
	r := config.SetupRouter(serverServices, handlers, middlewares.Authenticate(serverServices.Auth))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("an error occurred: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	log.Println("Shutting down gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Forced shutdown: %v", err)
	}

	log.Println("Shut down")
}
