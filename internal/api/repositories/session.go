package repositories

import (
	"context"

	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/types/subsonic"
	"gorm.io/gorm"
)

type SessionRepository struct {
	BaseRepository[models.Session]
}

func NewSessionRepository(Db *gorm.DB) *SessionRepository {
	return &SessionRepository{
		BaseRepository[models.Session]{Db: Db, Table: "sessions"},
	}
}

// FindCurrentByEntry find recent/active sessions for a given combo of user + track
func (r *SessionRepository) FindCurrentByEntry(ctx context.Context, entry subsonic.Entry) (*models.Session, error) {
	res, err := gorm.G[models.Session](r.Db).
		Where("username = ?", entry.Username).
		Where("player_id = ?", entry.PlayerID).
		Where("title = ?", entry.Title).
		Where("completed_at IS null").
		Where("EXTRACT(EPOCH FROM (NOW() - created_at)) <= ?", entry.Duration).First(ctx)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
