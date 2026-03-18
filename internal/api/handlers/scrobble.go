package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/common"
	"github.com/rangodisco/yhar/internal/api/pollers"
)

type ScrobbleHandler struct {
	subsonic pollers.PlayerPoller
}

func NewScrobbleHandler(subsonic pollers.PlayerPoller) *ScrobbleHandler {
	return &ScrobbleHandler{subsonic: subsonic}
}

func (h *ScrobbleHandler) ManualNowPlayingPoll(c *gin.Context) {
	ctx := c.Request.Context()

	err := h.subsonic.PollPlaying(ctx)
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, err, "Unable to poll")
		return
	}

	c.JSON(200, gin.H{
		"message": "polled subsonic successfully",
	})
	return
}
