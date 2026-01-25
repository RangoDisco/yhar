package repositories

import (
	"github.com/rangodisco/yhar/internal/api/config/database"
	"github.com/rangodisco/yhar/internal/api/models"
)

func FindActiveGenreByName(name string) (*models.Genre, error) {
	var g models.Genre
	err := database.GetDB().First(&g, "name = ?", name).Error
	if err != nil {
		return nil, err
	}

	return &g, nil
}

func CreateGenre(g *models.Genre) error {
	res := database.GetDB().Create(g)
	return res.Error
}
