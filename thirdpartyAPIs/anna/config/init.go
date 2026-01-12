package config

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/thirdpartyAPIs/anna/config/database"
)

func Init() *gin.Engine {
	err := database.SetupDatabase()
	if err != nil {
		log.Fatal(err)
	}

	r := SetupRouter()
	return r
}
