package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/config"
	ydb "github.com/rangodisco/yhar/internal/api/config/database"
	mdb "github.com/rangodisco/yhar/internal/metadata/config/database"
)

func init() {
	err := config.LoadEnv()
	if err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}
}

func main() {

	err := mdb.SetupDatabase()
	if err != nil {
		log.Fatalf("failed to init database: %v", err)
	}

	err = ydb.SetupDatabase()
	if err != nil {
		log.Fatalf("failed to init database: %v", err)
	}

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := config.SetupRouter()
	err = r.Run()
	if err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}
