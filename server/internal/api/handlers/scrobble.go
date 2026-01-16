package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/metadata/types/scrobble"
	"github.com/rangodisco/yhar/server/internal/services"
)

func GetNowPlaying(c *gin.Context) {
	var annaRes scrobble.InfoResponse
	subRes, err := services.GetNowPlaying()
	if err != nil {
		c.JSON(500, gin.H{})
	}

	if len(subRes.NowPlaying.Entry) == 0 {
		c.JSON(200, gin.H{})
	}

	for _, entry := range subRes.NowPlaying.Entry {
		data, err := services.GetTrackMetadata(&entry)
		if err != nil {
			// TODO: log
			continue
		}
		fmt.Println(data)
	}

	c.JSON(200, gin.H{
		"data": annaRes,
	})
}
