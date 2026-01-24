package services

import (
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/repositories"
)

func GetOrCreateGenre(name string) (*models.Genre, error) {
	existingGenre, err := repositories.FindActiveGenreByName(name)
	if err == nil {
		return existingGenre, nil
	}

	model := &models.Genre{
		Name: name,
	}
	err = repositories.CreateGenre(model)
	if err != nil {
		return nil, err
	}

	return model, nil
}
