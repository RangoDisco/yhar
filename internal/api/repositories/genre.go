package repositories

import (
	"context"

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

func (r *GenreRepository) FindActiveByName(ctx context.Context, name string) (*models.Genre, error) {
	var g models.Genre
	err := r.Db.WithContext(ctx).First(&g, "name = ?", name).Error
	if err != nil {
		return nil, err
	}

	return &g, nil
}

func (r *GenreRepository) CreateGenre(ctx context.Context, g *models.Genre) error {
	res := r.Db.WithContext(ctx).Create(g)
	return res.Error
}
