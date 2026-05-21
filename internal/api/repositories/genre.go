package repositories

import (
	"github.com/rangodisco/yhar/internal/api/models"
	"gorm.io/gorm"
)

type GenreRepository struct {
	BaseRepository[models.Genre]
}

func NewGenreRepository(Db *gorm.DB) *GenreRepository {
	return &GenreRepository{
		BaseRepository[models.Genre]{Db: Db, Table: "genres"},
	}
}
