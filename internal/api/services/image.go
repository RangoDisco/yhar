package services

import (
	"context"

	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/repositories"
)

type ImageService struct {
	repo *repositories.ImageRepository
}

func NewImageService(repo *repositories.ImageRepository) *ImageService {
	return &ImageService{repo: repo}
}

// GetOrCreate looks for the url in database, if it doesn't exist, creates, persists and returns the image
func (s *ImageService) GetOrCreate(ctx context.Context, url string) (*models.Image, error) {
	existingImage, err := s.repo.FindActiveImageByUrl(ctx, url)
	if err == nil && existingImage.Url != "" {
		return existingImage, nil
	}

	model := &models.Image{
		Url: url,
	}

	err = s.repo.PersistImage(ctx, model)
	if err != nil {
		return nil, err
	}
	return model, nil
}
