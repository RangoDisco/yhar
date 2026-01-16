package metadata

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/metadata/services"
	"github.com/rangodisco/yhar/internal/metadata/types/scrobble"
)

func GetTrackInfoByScrobble(c *gin.Context) {
	var body scrobble.InfoRequest
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	scrobbleInfo, err := services.GetInfoByScrobble(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{
			"scrobble": scrobbleInfo,
		},
	})
	return
}
