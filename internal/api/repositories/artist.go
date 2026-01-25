package repositories

import (
	"github.com/rangodisco/yhar/internal/api/config/database"
	"github.com/rangodisco/yhar/internal/api/models"
)

func FindActiveArtistByName(name string) (*models.Artist, error) {
	var a models.Artist

	err := database.GetDB().Where("name = ?", name).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, err
}

func PersistArtist(a *models.Artist) error {
	res := database.GetDB().Create(&a)
	return res.Error
}
