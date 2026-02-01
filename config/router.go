package config

import (
	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/handlers"
	"github.com/rangodisco/yhar/internal/api/middlewares"
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

	// AUTH
	auth := api.Group("/auth")
	auth.POST("/login", handlers.Login)

	protected := api.Group("/")
	protected.Use(middlewares.Authenticate())

	// THIRDPARTY
	navidrome := protected.Group("/navidrome")
	navidrome.GET("/getNowPlaying", middlewares.RequirePermissions([]string{"MANUAL_SCROBBLE"}), handlers.ManualNowPlayingPoll)

	// USER DATA
	user := protected.Group("/users/:userID")
	user.Use(middlewares.CheckUserPrivacy())

	// USER'S STATS
	userScrobbles := user.Group("/scrobbles")
	userScrobbles.GET("/top/artists", handlers.GetUserTopArtists)
	userScrobbles.GET("/top/albums", handlers.GetUserTopAlbums)
	userScrobbles.GET("/top/tracks", handlers.GetUserTopTracks)
}
