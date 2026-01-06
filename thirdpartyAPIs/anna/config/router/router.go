package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/thirdpartyAPIs/anna/internal/api/handlers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	handlers.RegisterTrackRoutes(r)

	return r
}
