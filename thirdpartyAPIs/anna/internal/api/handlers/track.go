package handlers

import "github.com/gin-gonic/gin"

func RegisterTrackRoutes(r *gin.Engine) {
	r.GET("/track-by-scrobble-details", GetTrackByScrobbleDetails)
}

func GetTrackByScrobbleDetails(c *gin.Context) {}
