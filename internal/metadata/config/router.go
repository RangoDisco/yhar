package config

import (
	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/handlers/metadata"
)

func SetupRouter() *gin.Engine {
	SetupLogger()
	r := gin.New()
	loadRoutes(r)

	// TOOD: middleware
	return r
}

func loadRoutes(r *gin.Engine) {
	api := r.Group("/api")
	tracks := api.Group("/tracks")
	tracks.POST("/by-scrobble", metadata.GetTrackInfoByScrobble)
}
