package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/rangodisco/yhar/internal/api/models"
	"gorm.io/gorm"
)

type ScrobbleRepository struct {
	Db *gorm.DB
}

func NewScrobbleRepository(Db *gorm.DB) *ScrobbleRepository {
	return &ScrobbleRepository{
		Db: Db,
	}
}

func (r *ScrobbleRepository) FindByTrackAndTimestamp(ctx context.Context, trackID, userID int64, timestamp time.Time) (*models.Scrobble, error) {
	var existingScrobble *models.Scrobble
	err := r.Db.WithContext(ctx).
		Where("track_id = ? AND user_id = ? AND timestamp = ?", trackID, userID, timestamp).
		Find(existingScrobble).Error
	if err != nil {
		return nil, err
	}

	return existingScrobble, nil

}

func (r *ScrobbleRepository) PersistScrobble(ctx context.Context, s *models.Scrobble) error {
	err := r.Db.WithContext(ctx).Create(&s).Error

	if err != nil {
		return fmt.Errorf("unable to persist scrobble: %w", err)
	}
	return nil
}
