package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/thirdpartyAPIs/anna/internal/api/handlers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// TOOD: middleware
	return r
}

func LoadRoutes(r *gin.Engine) {
	api := r.Group("/api")
	tracks := api.Group("/tracks")
	tracks.GET("/by-scrobble", handlers.GetTrackInfoByScrobble)
}
