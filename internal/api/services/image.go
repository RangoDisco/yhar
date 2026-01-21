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

	newImage, err := repositories.PersistImage(url)
	if err != nil {
		return nil, err
	}
	return newImage, nil
}
