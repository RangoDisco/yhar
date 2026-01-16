package repositories

import (
	"github.com/rangodisco/yhar/internal/metadata/config/database"
	"github.com/rangodisco/yhar/internal/metadata/models"
)

func FindAlbumById(id int64) (*[]models.Album, error) {
	var a []models.Album
	err := database.GetDB().Preload("Images").Preload("Artists.Images").Where("id = ?", id).Find(&a).Error
	if err != nil {
		return nil, err
	}

	return &a, nil
}
