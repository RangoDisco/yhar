package config

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rangodisco/yhar/thirdpartyAPIs/anna/config/database"
)

func Init() (*http.Server, error) {
	err := database.SetupDatabase()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to setup anna's database: %v", err))
	}

	r := SetupRouter()
	LoadRoutes(r)

	return &http.Server{
		Addr:         ":8081",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}, nil
}
