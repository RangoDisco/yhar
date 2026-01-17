package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/services"
	"github.com/rangodisco/yhar/internal/metadata/types/scrobble"
)

func ManualNowPlayingPoll(c *gin.Context) {
	var res []*scrobble.InfoResponse
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
		res = append(res, data)
	}

	c.JSON(200, gin.H{
		"data": res,
	})
}
