package config

import (
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	//err := database.SetupDatabase()
	//if err != nil {
	//	log.Fatal(err)
	//}

	r := SetupRouter()
	return r
}
