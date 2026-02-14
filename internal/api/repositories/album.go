package repositories

import (
	"context"

	"github.com/rangodisco/yhar/internal/api/models"
	"gorm.io/gorm"
)

type AlbumRepository struct {
	Db *gorm.DB
}

func NewAlbumRepository(Db *gorm.DB) *AlbumRepository {
	return &AlbumRepository{Db: Db}
}

func (r *AlbumRepository) FindActiveAlbumByTitle(ctx context.Context, title string) (*models.Album, error) {
	var a models.Album
	err := r.Db.WithContext(ctx).Preload("Artists.Images").Preload("Images").First(&a, "title = ?", title).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AlbumRepository) PersistAlbum(ctx context.Context, album *models.Album) error {
	res := r.Db.WithContext(ctx).Create(&album)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
