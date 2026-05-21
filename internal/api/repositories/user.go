package repositories

import (
	"github.com/rangodisco/yhar/internal/api/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	BaseRepository[models.User]
}

func NewUserRepository(Db *gorm.DB) *UserRepository {
	return &UserRepository{
		BaseRepository[models.User]{Db: Db, Table: "users"},
	}
}
