package main

import (
	"log"
	"net/http"
	"time"

	"github.com/rangodisco/yhar/config"
	anna "github.com/rangodisco/yhar/thirdpartyAPIs/anna/config"
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
	annaServer := &http.Server{
		Addr:         ":8081",
		Handler:      anna.Init(),
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
	}

	g.Go(func() error {
		return annaServer.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
