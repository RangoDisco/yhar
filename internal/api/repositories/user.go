package repositories

import (
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/types/filters"
	"gorm.io/gorm"
)

type UserRepository struct {
	Db *gorm.DB
}

func NewUserRepository(Db *gorm.DB) *UserRepository {
	return &UserRepository{Db: Db}
}

func (r *UserRepository) FindActiveUserByFilters(filters []filters.QueryFilter) (*models.User, error) {
	var u models.User
	query := r.Db.Preload("Role.Permissions")

	for _, filter := range filters {
		query.Where(filter.Key+" = ?", filter.Value)
	}

	err := query.First(&u).Error
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) PersistUser(user *models.User) error {
	res := r.Db.Create(user)
	return res.Error
}
