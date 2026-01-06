package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rangodisco/yhar/config"
	anna "github.com/rangodisco/yhar/thirdpartyAPIs/anna/config"
	annaDB "github.com/rangodisco/yhar/thirdpartyAPIs/anna/config/database"
	annaRouter "github.com/rangodisco/yhar/thirdpartyAPIs/anna/config/router"
	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

func init() {
	err := config.LoadEnv()
	if err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}
}

func main() {
	err := annaDB.SetupDatabase()
	if err != nil {
		log.Fatalf("failed to setup anna's database: %v", err)
	}

	r := annaRouter.SetupRouter()
	annaRouter.LoadRoutes(r)

	err = r.Run()
	if err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

	annaServer := &http.Server{
		Addr:         ":8081",
		Handler:      anna.Router(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
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
