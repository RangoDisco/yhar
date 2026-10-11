package repositories

import (
	"context"
	"fmt"

	"github.com/rangodisco/yhar/api/models"
	"gorm.io/gorm"
)

type AlbumRepository struct {
	BaseRepository[models.Album]
}

func NewAlbumRepository(Db *gorm.DB) *AlbumRepository {
	return &AlbumRepository{
		BaseRepository: BaseRepository[models.Album]{Db: Db, Table: "albums"},
	}
}

func (r *AlbumRepository) FindOneByFilter(ctx context.Context, filters []QueryFilter, preloads ...string) (*models.Album, error) {
	var res *models.Album

	query := r.preload(r.Db.WithContext(ctx).Session(&gorm.Session{}), preloads).
		Table("albums a").
		Joins("INNER JOIN artist_albums al ON al.album_id = a.id").
		Joins("INNER JOIN artists ar ON al.artist_id = ar.id")

	for _, f := range filters {
		if f.Value == nil {
			query = query.Where(fmt.Sprintf("%s IS null", f.Key))
		} else {
			query = query.Where(fmt.Sprintf("%s = ?", f.Key), f.Value)
		}
	}

	err := query.First(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}
