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

func (r *ImageRepository) FindActiveImageByUrl(ctx context.Context, url string) (*models.Image, error) {
	var i models.Image

	err := r.Db.WithContext(ctx).Where("url = ?", url).First(&i).Error
	if err != nil {
		return nil, err
	}
	return &i, err
}

func (r *ImageRepository) PersistImage(ctx context.Context, img *models.Image) error {

	res := r.Db.WithContext(ctx).Create(img)
	return res.Error
}
