package repositories

import (
	"github.com/rangodisco/yhar/internal/api/models"
	"gorm.io/gorm"
)

type IAlbumRepository interface {
	FindActiveAlbumByTitle(title string) (album *models.Album, err error)
	PersistAlbum(album *models.Album) error
}

type AlbumRepository struct {
	Db *gorm.DB
}

func NewAlbumRepository(Db *gorm.DB) IAlbumRepository {
	return &AlbumRepository{Db: Db}
}

func (r *AlbumRepository) FindActiveAlbumByTitle(title string) (*models.Album, error) {
	var a models.Album
	err := r.Db.Preload("Artists.Images").Preload("Images").First(&a, "title = ?", title).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AlbumRepository) PersistAlbum(album *models.Album) error {
	res := r.Db.Create(&album)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
