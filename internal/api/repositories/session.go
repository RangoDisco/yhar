package repositories

import (
	"context"
	"fmt"

	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/types/filters"
	"gorm.io/gorm"
)

type SessionRepository struct {
	Db *gorm.DB
}

func NewSessionRepository(Db *gorm.DB) *SessionRepository {
	return &SessionRepository{
		Db: Db,
	}
}

func (r *SessionRepository) FindByFilters(ctx context.Context, filters []filters.QueryFilter) (*models.Session, error) {
	var s models.Session
	query := r.Db.WithContext(ctx).
		Where("completed_at IS null")

	for _, filter := range filters {
		query = query.Where(filter.Key+" = ?", filter.Value)
	}

	err := query.First(&s).Error
	if err != nil {
		return nil, fmt.Errorf("unable to find session: %w", err)
	}

	return &s, nil
}

func (r *SessionRepository) Persist(ctx context.Context, s *models.Session) error {
	err := r.Db.WithContext(ctx).Create(&s).Error

	if err != nil {
		return fmt.Errorf("unable to persist session: %w", err)
	}

	return nil
}

func (r *SessionRepository) Update(ctx context.Context, s *models.Session) error {
	err := r.Db.WithContext(ctx).Save(&s).Error

	if err != nil {
		return fmt.Errorf("unable to update session: %w", err)
	}

	return nil
}
