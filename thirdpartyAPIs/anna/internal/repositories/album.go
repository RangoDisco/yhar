package repositories

import (
	"github.com/rangodisco/yhar/thirdpartyAPIs/anna/config/database"
	"github.com/rangodisco/yhar/thirdpartyAPIs/anna/internal/models"
)

func FindAlbumById(id int64) (*[]models.Album, error) {
	var a []models.Album
	err := database.GetDB().Preload("artists").Where("id = ?", id).Find(&a).Error
	if err != nil {
		return nil, err
	}

	return &a, nil
}
