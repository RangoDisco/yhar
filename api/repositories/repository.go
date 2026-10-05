package repositories

import (
	"context"

	"gorm.io/gorm"
)

type QueryFilter struct {
	Key   string
	Value any
}

type Scope = func(*gorm.DB) *gorm.DB

type Repository[T any] interface {
	Persist(ctx context.Context, model *T) error
	FindOneByID(ctx context.Context, id int64, preloads ...string) (*T, error)
	FindOneBy(ctx context.Context, filters []QueryFilter, preloads ...string) (*T, error)
	FindBy(ctx context.Context, filters []QueryFilter, preloads ...string) ([]T, error)
	Update(ctx context.Context, id int64, fields map[string]interface{}) error
	Delete(ctx context.Context, id int64) error
}
