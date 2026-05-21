package handlers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/common"
	"github.com/rangodisco/yhar/internal/api/dto"
	"github.com/rangodisco/yhar/internal/api/services"
	"gorm.io/gorm"
)

type ArtistHandler struct {
	service *services.ArtistService
}

func NewArtistHandler(service *services.ArtistService) *ArtistHandler {
	return &ArtistHandler{service: service}
}

// Update partially updates an artist with given fields
func (h *ArtistHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	var input dto.UpdateArtistInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		common.RespondWithError(c, 400, err, "Invalid body")
		return
	}

	iID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		common.RespondWithError(c, 400, err, "Invalid artist ID")
		return
	}

	artist, err := h.service.GetByID(ctx, iID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			common.RespondWithError(c, 404, err, "Artist not found")
		} else {
			common.RespondWithError(c, 500, err, "Unable to fetch artist")
		}
		return
	}

	updated, err := h.service.Update(ctx, artist, input)
	if err != nil {
		common.RespondWithError(c, 500, err, "Unable to update artist")
		return
	}

	common.RespondWithData(c, 200, updated)
}
