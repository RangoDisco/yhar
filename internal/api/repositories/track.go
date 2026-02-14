package repositories

import (
	"context"

	"github.com/rangodisco/yhar/internal/api/models"
	"gorm.io/gorm"
)

type TrackRepository struct {
	Db *gorm.DB
}

func NewTrackRepository(Db *gorm.DB) *TrackRepository {
	return &TrackRepository{Db: Db}
}

func (r *TrackRepository) FindActiveByTitle(ctx context.Context, title string) (*models.Track, error) {
	var t models.Track

	// TODO: handle multiple track with same name (check for albums/artists)
	err := r.Db.WithContext(ctx).Preload("Artists.Picture").Preload("Album.Picture").Where("title = ?", title).First(&t).Error
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *TrackRepository) PersistTrack(ctx context.Context, track *models.Track) error {
	res := r.Db.WithContext(ctx).Create(&track)

	return res.Error
}
