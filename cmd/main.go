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

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/config"
	serverConfig "github.com/rangodisco/yhar/internal/api/config"
	ydb "github.com/rangodisco/yhar/internal/api/config/database"
)

func init() {
	err := config.LoadEnv()
	if err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}

	if os.Getenv("JWT_SECRET") == "" {
		log.Fatalf("JWT_SECRET environment variable not set")
	}
}

func main() {
	yDb, err := ydb.SetupDatabase()
	if err != nil {
		log.Fatalf("failed to init database: %v", err)
	}

	switch os.Getenv("GIN_MODE") {
	case gin.DebugMode:
		gin.SetMode(gin.DebugMode)
		break
	case gin.TestMode:
		gin.SetMode(gin.TestMode)
		break
	case gin.ReleaseMode:
	default:
		gin.SetMode(gin.ReleaseMode)
		break
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	serverRepos, serverServices, handlers, pollers, importers := serverConfig.AutoWire(yDb)

	// Start pollers (subsonic only for now)
	if pollers.Subsonic != nil {
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

	// Start router
	r := config.SetupRouter(serverRepos, serverServices, handlers)

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
