package services

import (
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/repositories"
)

type ImageService struct {
	iRepo *repositories.ImageRepository
}

func NewImageService(iRepo *repositories.ImageRepository) *ImageService {
	return &ImageService{iRepo: iRepo}
}

// GetOrCreateImage looks for the url in database, if it doesn't exist, creates and returns it
func (s *ImageService) GetOrCreateImage(url string) (*models.Image, error) {
	existingImage, err := s.iRepo.FindActiveImageByUrl(url)
	if err == nil && existingImage.Url != "" {
		return existingImage, nil
	}

	model := buildImageModel(url)

	err = s.iRepo.PersistImage(model)
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
