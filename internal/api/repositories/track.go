package repositories

import (
	"context"
	"fmt"

	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/types/filters"
	"gorm.io/gorm"
)

type TrackRepository struct {
	Db *gorm.DB
}

func NewTrackRepository(Db *gorm.DB) *TrackRepository {
	return &TrackRepository{Db: Db}
}

func (r *TrackRepository) FindActiveByFilter(ctx context.Context, filters []filters.QueryFilter) (*models.Track, error) {
	var t models.Track

	if len(filters) == 0 {
		return nil, fmt.Errorf("filters are empty")
	}

	query := r.Db.WithContext(ctx).Preload("Artists.Picture").Preload("Album.Picture")

	for _, f := range filters {
		query = query.Where(fmt.Sprintf("%s = ?", f.Key), f.Value)
	}

	err := query.First(&t).Error
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *TrackRepository) Persist(ctx context.Context, track *models.Track) error {
	res := r.Db.WithContext(ctx).Create(&track)

	return res.Error
}
