package services

import (
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/repositories"
)

// GetOrCreateImage looks for the url in database, if it doesn't exist, creates and returns it
func GetOrCreateImage(url string) (*models.Image, error) {
	existingImage, err := repositories.FindActiveImageByUrl(url)
	if err == nil && existingImage.Url != "" {
		return existingImage, nil
	}

	model := buildImageModel(url)

	err = repositories.PersistImage(model)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func buildImageModel(url string) *models.Image {
	return &models.Image{
		Url: url,
	}
}
