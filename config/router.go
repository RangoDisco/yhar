package config

import (
	"github.com/gin-gonic/gin"
	serverConfig "github.com/rangodisco/yhar/internal/api/config"
	"github.com/rangodisco/yhar/internal/api/middlewares"
)

func SetupRouter(
	repo *serverConfig.Repositories,
	s *serverConfig.Services,
	h *serverConfig.Handlers,
) *gin.Engine {
	SetupLogger()
	r := gin.New()

	loadRoutes(r, repo, s, h)

	return r
}

func loadRoutes(r *gin.Engine, repo *serverConfig.Repositories, s *serverConfig.Services, h *serverConfig.Handlers) {
	api := r.Group("/api")

	// AUTH
	auth := api.Group("/auth")
	auth.POST("/login", h.Auth.Login)

	protected := api.Group("/")
	protected.Use(middlewares.Authenticate(s.Auth))

	// THIRDPARTY
	subsonic := protected.Group("/subsonic")
	subsonic.GET("/getNowPlaying", middlewares.RequirePermissions([]string{"MANUAL_SCROBBLE"}), h.Scrobble.ManualNowPlayingPoll)

	// USER DATA
	user := protected.Group("/users/:userID")
	user.Use(middlewares.CheckUserPrivacy(repo.User))

	user.GET("", h.User.GetUser)

	// USER'S STATS
	userScrobbles := user.Group("/scrobbles")
	userScrobbles.GET("/history", h.Scrobble.GetUserHistory)
	userScrobbles.GET("/top/artists", h.Scrobble.GetUserTopArtists)
	userScrobbles.GET("/top/albums", h.Scrobble.GetUserTopAlbums)
	userScrobbles.GET("/top/tracks", h.Scrobble.GetUserTopTracks)
}
