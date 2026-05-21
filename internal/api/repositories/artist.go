package repositories

import (
	"github.com/rangodisco/yhar/internal/api/models"
	"gorm.io/gorm"
)

type ArtistRepository struct {
	BaseRepository[models.Artist]
}

func NewArtistRepository(Db *gorm.DB) *ArtistRepository {
	return &ArtistRepository{
		BaseRepository[models.Artist]{Db: Db, Table: "artists"},
	}
}
