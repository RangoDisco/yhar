package repositories

import (
	"github.com/rangodisco/yhar/internal/metadata/models"
	"gorm.io/gorm"
)

type AlbumRepository struct {
	Db *gorm.DB
}

func NewAlbumRepository(db *gorm.DB) *AlbumRepository {
	return &AlbumRepository{
		Db: db,
	}
}

func (r *AlbumRepository) FindAlbumById(id int64) (*[]models.Album, error) {
	var a []models.Album
	err := r.Db.Preload("Images").Preload("Artists.Images").Where("id = ?", id).Find(&a).Error
	if err != nil {
		return nil, err
	}

	return &a, nil
}
