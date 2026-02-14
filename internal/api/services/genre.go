package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/repositories"
	"gorm.io/gorm"
)

type GenreService struct {
	repo *repositories.GenreRepository
}

func NewGenreService(repo *repositories.GenreRepository) *GenreService {
	return &GenreService{repo: repo}
}

// GetOrCreateGenre tries to find a genre by its name, and creates if it doesn't already exist
func (s *GenreService) GetOrCreateGenre(ctx context.Context, name string) (*models.Genre, error) {
	existingGenre, err := s.repo.FindActiveByName(ctx, name)
	if err == nil {
		return existingGenre, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to find genre: %w", err)
	}

	model := &models.Genre{
		Name: name,
	}
	err = s.repo.CreateGenre(ctx, model)
	if err != nil {
		return nil, err
	}

	return model, nil
}
