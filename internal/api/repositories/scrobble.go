package repositories

import (
	"context"
	"fmt"

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

func (r *ScrobbleRepository) PersistScrobble(ctx context.Context, s *models.Scrobble) error {
	err := r.Db.WithContext(ctx).Create(&s).Error

	if err != nil {
		return fmt.Errorf("unable to persist scrobble: %w", err)
	}
	return nil
}
