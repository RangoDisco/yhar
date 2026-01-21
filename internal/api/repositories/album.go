package repositories

import (
	"github.com/rangodisco/yhar/internal/api/config/database"
	"github.com/rangodisco/yhar/internal/api/models"
)

func FindActiveAlbumByTitle(title string) (*models.Album, error) {
	var a models.Album
	err := database.GetDB().Preload("Artists.Images").Preload("Images").First(&a, "title = ?", title).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func PersistAlbum(album *models.Album) (*models.Album, error) {
	res := database.GetDB().Create(&album)
	if res.Error != nil {
		return nil, res.Error
	}
	return album, nil
}
