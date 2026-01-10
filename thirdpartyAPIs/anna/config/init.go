package config

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rangodisco/yhar/thirdpartyAPIs/anna/config/database"
	"github.com/rangodisco/yhar/thirdpartyAPIs/anna/config/router"
)

func Init() (*http.Server, error) {
	err := database.SetupDatabase()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to setup anna's database: %v", err))
	}

	r := router.SetupRouter()
	router.LoadRoutes(r)

	err = r.Run()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to start anna's router: %v", err))
	}

	SetupLogger()

	return &http.Server{
		Addr:         ":8081",
		Handler:      Router(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}, nil
}
