package repositories

import (
	"github.com/rangodisco/yhar/internal/api/models"
	"gorm.io/gorm"
)

type TrackRepository struct {
	BaseRepository[models.Track]
}

func NewTrackRepository(Db *gorm.DB) *TrackRepository {
	return &TrackRepository{
		BaseRepository[models.Track]{Db: Db, Table: "tracks"},
	}
}
