package repositories

import (
	"github.com/rangodisco/yhar/internal/api/models"
	"gorm.io/gorm"
)

type ScrobbleRepository struct {
	BaseRepository[models.Scrobble]
}

func NewScrobbleRepository(Db *gorm.DB) *ScrobbleRepository {
	return &ScrobbleRepository{
		BaseRepository[models.Scrobble]{Db: Db, Table: "scrobbles"},
	}
}

// TODO: remove
//func (r *ScrobbleRepository) FindByTrackAndTimestamp(ctx context.Context, trackID, userID int64, timestamp time.Time) (*models.Scrobble, error) {
//	var existingScrobble *models.Scrobble
//	err := r.Db.WithContext(ctx).
//		Where("track_id = ? AND user_id = ? AND scrobbled_at = ?", trackID, userID, timestamp).
//		First(&existingScrobble).Error
//	if err != nil {
//		return nil, err
//	}
//
//	return existingScrobble, nil
//
//}
