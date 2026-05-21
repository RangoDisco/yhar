package repositories

import (
	"github.com/rangodisco/yhar/internal/api/models"
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
