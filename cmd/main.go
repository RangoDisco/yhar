package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/config"
	serverConfig "github.com/rangodisco/yhar/internal/api/config"
	ydb "github.com/rangodisco/yhar/internal/api/config/database"
	metaConfig "github.com/rangodisco/yhar/internal/metadata/config"
	mdb "github.com/rangodisco/yhar/internal/metadata/config/database"
)

func init() {
	err := config.LoadEnv()
	if err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}
}

func main() {
	mDb, err := mdb.SetupDatabase()
	if err != nil {
		log.Fatalf("failed to init database: %v", err)
	}

	yDb, err := ydb.SetupDatabase()
	if err != nil {
		log.Fatalf("failed to init database: %v", err)
	}

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// TODO: Pls split it in two different git repository
	metaServices := metaConfig.AutoWire(mDb)
	serverRepos, serverServices, handlers := serverConfig.AutoWire(yDb, metaServices)

	r := config.SetupRouter(serverRepos, serverServices, handlers)
	err = r.Run()
	if err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}
