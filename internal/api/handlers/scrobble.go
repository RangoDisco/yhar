package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/services"
)

type ScrobbleHandler struct {
	service  *services.ScrobbleService
	subsonic *services.SubsonicService
}

func NewScrobbleHandler(service *services.ScrobbleService, subsonic *services.SubsonicService) *ScrobbleHandler {
	return &ScrobbleHandler{service: service, subsonic: subsonic}
}

// TODO: handle polling and what is considered a real scrobble
func (h *ScrobbleHandler) ManualNowPlayingPoll(c *gin.Context) {
	ctx := c.Request.Context()
	var scrobbles []*models.Scrobble
	user, exists := c.Get("user")
	if !exists {
		c.JSON(500, gin.H{})
		return
	}
	fmt.Println(user)

	subRes, err := h.subsonic.GetNowPlaying()
	if err != nil {
		c.JSON(500, gin.H{})
		return
	}

	if len(subRes.NowPlaying.Entry) == 0 {
		c.JSON(200, gin.H{})
		return
	}

	for _, entry := range subRes.NowPlaying.Entry {
		// Check if track already exists in service's db
		res, err := h.service.HandleNewScrobble(ctx, entry)
		if err != nil {
			continue
		}
		scrobbles = append(scrobbles, res)
	}

	c.JSON(200, gin.H{
		"data": scrobbles,
	})
}
