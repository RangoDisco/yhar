package repositories

import (
	"github.com/rangodisco/yhar/internal/api/config/database"
	"github.com/rangodisco/yhar/internal/api/models"
)

func FindActiveImageByUrl(url string) (*models.Image, error) {
	var i models.Image

	err := database.GetDB().Where("url = ?", url).First(&i).Error
	if err != nil {
		return nil, err
	}
	return &i, err
}

func PersistImage(url string) (*models.Image, error) {
	img := &models.Image{
		Url: url,
	}
	err := database.GetDB().Create(&img).Error
	if err != nil {
		return nil, err
	}
	return img, err
}
