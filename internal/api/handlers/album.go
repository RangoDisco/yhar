package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/common"
	"github.com/rangodisco/yhar/internal/api/dto"
	"github.com/rangodisco/yhar/internal/api/services"
	"github.com/rangodisco/yhar/internal/api/types/filters"
	"gorm.io/gorm"
)

type AlbumHandler struct {
	service *services.AlbumService
}

func NewAlbumHandler(service *services.AlbumService) *AlbumHandler {
	return &AlbumHandler{service: service}
}

func (h *AlbumHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	var input dto.UpdateAlbumInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		common.RespondWithError(c, 400, err, "Invalid body")
		return
	}

	album, err := h.service.Get(ctx, []filters.QueryFilter{{Key: "id", Value: id}})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			common.RespondWithError(c, 404, err, "Album not found")
		} else {
			common.RespondWithError(c, 500, err, "Unable to fetch album")
		}
		return
	}

	updated, err := h.service.Update(ctx, album, input)
	if err != nil {
		common.RespondWithError(c, 500, err, "Unable to update album")
		return
	}

	common.RespondWithData(c, 200, updated)
}
