package repositories

import (
	"context"

	"github.com/rangodisco/yhar/internal/api/models"
	"gorm.io/gorm"
)

type ImageRepository struct {
	Db *gorm.DB
}

func NewImageRepository(Db *gorm.DB) *ImageRepository {
	return &ImageRepository{Db: Db}
}

func (r *ImageRepository) Persist(ctx context.Context, img *models.Image) error {

	res := r.Db.WithContext(ctx).Create(img)
	return res.Error
}
