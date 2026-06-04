package handlers

import (
	"fmt"
	"log"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/common"
	"github.com/rangodisco/yhar/internal/api/services"
)

type ImageHandler struct {
	service *services.ImageService
}

func NewImageHandler(service *services.ImageService) *ImageHandler {
	return &ImageHandler{service: service}
}

func (h *ImageHandler) Upload(c *gin.Context) {
	ctx := c.Request.Context()
	file, err := c.FormFile("image")
	if err != nil {
		common.RespondWithError(c, 400, err, "No file provided")
		return
	}

	src, err := file.Open()
	if err != nil {
		common.RespondWithError(c, 400, err, "Unable to open provided file")
		return
	}

	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			log.Printf("unable to close file: %v", err)
		}
	}(src)

	// Pre validate size so we can safely assign the whole file size later on
	if file.Size > services.DefaultConfig.MaxSize {
		common.RespondWithError(c, 400, fmt.Errorf("maximum size is: %d, got %d", services.DefaultConfig.MaxSize, file.Size), "File is too big")
		return
	}

	buf := make([]byte, file.Size)
	if _, err = src.Read(buf); err != nil {
		common.RespondWithError(c, 400, err, "Unable to read provided file")
		return
	}

	err = h.service.Validate(file, buf)
	if err != nil {
		common.RespondWithError(c, 400, err, "File is invalid")
		return
	}

	filename, err := h.service.SaveLocally(ctx, buf)
	if err != nil {
		common.RespondWithError(c, 500, err, "Unable to save image locally")
		return
	}

	image, err := h.service.Persist(ctx, filename)
	if err != nil {
		common.RespondWithError(c, 500, err, "Unable to persist image")
		return
	}

	common.RespondWithData(c, 200, image)
}
