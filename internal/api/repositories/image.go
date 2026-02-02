package repositories

import (
	"github.com/rangodisco/yhar/internal/api/config/database"
	"github.com/rangodisco/yhar/internal/api/models"
	"gorm.io/gorm"
)

type ImageRepository struct {
	Db *gorm.DB
}

func NewImageRepository(Db *gorm.DB) *ImageRepository {
	return &ImageRepository{Db: Db}
}

func (r *ImageRepository) FindActiveImageByUrl(url string) (*models.Image, error) {
	var i models.Image

	err := database.GetDB().Where("url = ?", url).First(&i).Error
	if err != nil {
		return nil, err
	}
	return &i, err
}

func (r *ImageRepository) PersistImage(img *models.Image) error {

	res := database.GetDB().Create(img)
	return res.Error
}
