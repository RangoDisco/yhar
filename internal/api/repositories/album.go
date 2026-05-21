package repositories

import (
	"github.com/rangodisco/yhar/internal/api/models"
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
