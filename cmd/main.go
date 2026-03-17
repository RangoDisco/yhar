package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

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
}

func main() {
	yDb, err := ydb.SetupDatabase()
	if err != nil {
		log.Fatalf("failed to init database: %v", err)
	}

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
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

	go func() {
		if err := r.Run(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	log.Println("Shutting down gracefully...")
	log.Println("Shut down")
}
