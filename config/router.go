package config

import (
	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/middlewares"
)

func SetupRouter(repo *Repositories, s *Services, h *Handlers) *gin.Engine {
	SetupLogger()
	r := gin.New()

	loadRoutes(r, repo, s, h)

	return r
}

func loadRoutes(r *gin.Engine, repo *Repositories, s *Services, h *Handlers) {
	api := r.Group("/api")

	// AUTH
	auth := api.Group("/auth")
	auth.POST("/login", h.Auth.Login)

	protected := api.Group("/")
	protected.Use(middlewares.Authenticate(s.Auth))

	// THIRDPARTY
	navidrome := protected.Group("/navidrome")
	navidrome.GET("/getNowPlaying", middlewares.RequirePermissions([]string{"MANUAL_SCROBBLE"}), h.Scrobble.ManualNowPlayingPoll)

	// USER DATA
	user := protected.Group("/users/:userID")
	user.Use(middlewares.CheckUserPrivacy(repo.User))

	// USER'S STATS
	userScrobbles := user.Group("/scrobbles")
	userScrobbles.GET("/top/artists", h.Scrobble.GetUserTopArtists)
	userScrobbles.GET("/top/albums", h.Scrobble.GetUserTopAlbums)
	userScrobbles.GET("/top/tracks", h.Scrobble.GetUserTopTracks)
}
