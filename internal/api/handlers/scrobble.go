package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/common"
	"github.com/rangodisco/yhar/internal/api/pollers"
	"github.com/rangodisco/yhar/internal/api/services"
	"gorm.io/gorm"
)

type ScrobbleHandler struct {
	subsonic pollers.PlayerPoller
	service  *services.ScrobbleService
}

func NewScrobbleHandler(subsonic pollers.PlayerPoller, service *services.ScrobbleService) *ScrobbleHandler {
	return &ScrobbleHandler{subsonic: subsonic, service: service}
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

func (h *ScrobbleHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	currentUser, err := common.GetUserFromContext(c)
	if err != nil {
		common.RespondWithError(c, 401, errors.New("user is not authenticated"), "Unable to fetch user")
		return
	}

	iID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		common.RespondWithError(c, 400, err, "Invalid scrobble ID")
		return
	}

	scrobble, err := h.service.GetByID(ctx, iID, currentUser)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			common.RespondWithError(c, 404, err, "Scrobble not found")
		} else {
			common.RespondWithError(c, 500, err, "Unable to fetch scrobble")
		}
		return
	}

	err = h.service.Delete(ctx, scrobble)
	if err != nil {
		common.RespondWithError(c, 500, err, "Unable to delete scrobble")
		return
	}

	common.RespondWithData(c, 200, gin.H{
		"success": true,
	})
}
