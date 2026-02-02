package repositories

import (
	"github.com/rangodisco/yhar/internal/api/models"
	"gorm.io/gorm"
)

type GenreRepository struct {
	Db *gorm.DB
}

func NewGenreRepository(Db *gorm.DB) *GenreRepository {
	return &GenreRepository{
		Db: Db,
	}
}

func (r *GenreRepository) FindActiveGenreByName(name string) (*models.Genre, error) {
	var g models.Genre
	err := r.Db.First(&g, "name = ?", name).Error
	if err != nil {
		return nil, err
	}

	return &g, nil
}

func (r *GenreRepository) CreateGenre(g *models.Genre) error {
	res := r.Db.Create(g)
	return res.Error
}
