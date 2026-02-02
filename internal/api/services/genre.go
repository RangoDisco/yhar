package services

import (
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/repositories"
)

type GenreService struct {
	gRepo *repositories.GenreRepository
}

func NewGenreService(gRepo *repositories.GenreRepository) *GenreService {
	return &GenreService{gRepo: gRepo}
}

func (s *GenreService) GetOrCreateGenre(name string) (*models.Genre, error) {
	existingGenre, err := s.gRepo.FindActiveGenreByName(name)
	if err == nil {
		return existingGenre, nil
	}

	model := &models.Genre{
		Name: name,
	}
	err = s.gRepo.CreateGenre(model)
	if err != nil {
		return nil, err
	}

	return model, nil
}
