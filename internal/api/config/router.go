package config

import (
	"time"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/middlewares"
)

func SetupRouter(
	s *Services,
	h *Handlers,
	authMiddleware gin.HandlerFunc,
) *gin.Engine {
	r := gin.New()

	api := r.Group("/api")
	api.Use(sentrygin.New(sentrygin.Options{Repanic: true}))
	api.Use(middlewares.LoggerMiddleware())
	api.Use(cors.Default())
	api.Use(middlewares.SecurityMiddleware())
	api.Use(middlewares.RateLimiter())
	api.Use(middlewares.TimeoutMiddleware(10 * time.Second))

	r.Static("/images", "images")

	// AUTH
	auth := api.Group("/auth")
	auth.POST("/login", h.Auth.Login)
	auth.POST("/refresh", h.Auth.Refresh)

	protected := api.Group("/")
	protected.Use(authMiddleware)

	// Image
	protected.POST("/images", middlewares.RequirePermissions([]string{"IMAGE_UPLOAD"}), h.Image.Upload)

	// Artists
	protected.PATCH("/artists/:id", middlewares.RequirePermissions([]string{"UPDATE_ARTIST"}), h.Artist.Update)

	// Albums
	protected.PATCH("/albums/:id", middlewares.RequirePermissions([]string{"UPDATE_ALBUM"}), h.Album.Update)

	// THIRDPARTY
	subsonic := protected.Group("/subsonic")
	subsonic.GET("/getNowPlaying", middlewares.RequirePermissions([]string{"MANUAL_SCROBBLE"}), h.Scrobble.ManualNowPlayingPoll)

	// CRUD
	protected.DELETE("/scrobbles/:id", h.Scrobble.Delete)

	// USER DATA
	user := protected.Group("/users/:userID")
	user.Use(middlewares.CheckUserPrivacy(s.User))

	user.GET("", h.User.GetUser)

	// USER'S STATS
	userScrobbles := user.Group("/scrobbles")
	userScrobbles.GET("/history", h.ScrobbleStats.GetUserHistory)
	userScrobbles.GET("/top/artists", h.ScrobbleStats.GetUserTopArtists)
	userScrobbles.GET("/top/albums", h.ScrobbleStats.GetUserTopAlbums)
	userScrobbles.GET("/top/tracks", h.ScrobbleStats.GetUserTopTracks)

	return r
}
