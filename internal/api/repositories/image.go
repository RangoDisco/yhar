package repositories

import (
	"github.com/rangodisco/yhar/internal/api/models"
	"gorm.io/gorm"
)

type ImageRepository struct {
	BaseRepository[models.Image]
}

func NewImageRepository(Db *gorm.DB) *ImageRepository {
	return &ImageRepository{
		BaseRepository[models.Image]{Db: Db, Table: "images"},
	}
}
