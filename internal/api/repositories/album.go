package repositories

import (
	"context"
	"fmt"

	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/types/filters"
	"gorm.io/gorm"
)

type AlbumRepository struct {
	Db *gorm.DB
}

func NewAlbumRepository(Db *gorm.DB) *AlbumRepository {
	return &AlbumRepository{Db: Db}
}

func (r *AlbumRepository) FindActiveByFilters(ctx context.Context, filters []filters.QueryFilter) (*models.Album, error) {
	if len(filters) == 0 {
		return nil, fmt.Errorf("no filters provided")
	}
	var a models.Album

	query := r.Db.WithContext(ctx).Preload("Artists.Images").Preload("Images")

	for _, f := range filters {
		query.Where(fmt.Sprintf("%s = ?", f.Key), f.Value)
	}

	err := query.First(&a).Error
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
