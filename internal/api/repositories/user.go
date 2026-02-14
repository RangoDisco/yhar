package repositories

import (
	"context"

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

func (r *UserRepository) FindActiveByFilters(ctx context.Context, filters []filters.QueryFilter) (*models.User, error) {
	var u models.User
	query := r.Db.WithContext(ctx).Preload("Role.Permissions")

	for _, filter := range filters {
		query.Where(filter.Key+" = ?", filter.Value)
	}

	err := query.First(&u).Error
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) Persist(ctx context.Context, user *models.User) error {
	res := r.Db.WithContext(ctx).Create(user)
	return res.Error
}
