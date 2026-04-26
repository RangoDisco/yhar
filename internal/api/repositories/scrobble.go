package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/types/filters"
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
		Where("track_id = ? AND user_id = ? AND scrobbled_at = ?", trackID, userID, timestamp).
		First(&existingScrobble).Error
	if err != nil {
		return nil, err
	}

	return existingScrobble, nil

}

func (r *ScrobbleRepository) FindActiveByFilters(ctx context.Context, filters []filters.QueryFilter) (*models.Scrobble, error) {
	var s models.Scrobble

	if len(filters) == 0 {
		return nil, fmt.Errorf("filters are empty")
	}

	query := r.Db.WithContext(ctx)

	for _, f := range filters {
		query = query.Where(fmt.Sprintf("%s = ?", f.Key), f.Value)
	}

	err := query.First(&s).Error
	if err != nil {
		return nil, err
	}

	return &s, err
}

func (r *ScrobbleRepository) Delete(ctx context.Context, scrobble *models.Scrobble) error {
	err := r.Db.WithContext(ctx).Delete(&scrobble).Error
	if err != nil {
		return fmt.Errorf("unable to delete scrobble: %w", err)
	}
	return nil
}

func (r *ScrobbleRepository) PersistScrobble(ctx context.Context, s *models.Scrobble) error {
	err := r.Db.WithContext(ctx).Create(&s).Error

	if err != nil {
		return fmt.Errorf("unable to persist scrobble: %w", err)
	}
	return nil
}
