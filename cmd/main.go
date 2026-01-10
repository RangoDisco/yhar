package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rangodisco/yhar/config"
	anna "github.com/rangodisco/yhar/thirdpartyAPIs/anna/config"
)

func init() {
	err := config.LoadEnv()
	if err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}
}

func main() {

	annaServer, err := anna.Init()
	if err != nil {
		log.Fatalf(err.Error())
	}

	go func() {
		err := annaServer.ListenAndServe()
		if err != nil {
			log.Fatalf("failed to start anna server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := annaServer.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err.Error())
	}
}
