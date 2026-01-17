package config

import (
	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/handlers"
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
	navidrome := api.Group("/navidrome")
	navidrome.GET("/getNowPlaying", handlers.GetNowPlaying)
}
