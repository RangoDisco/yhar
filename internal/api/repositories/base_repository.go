package repositories

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type BaseRepository[T any] struct {
	Db    *gorm.DB
	Table string
}

func (r *BaseRepository[T]) preload(db *gorm.DB, preloads []string) *gorm.DB {
	for _, p := range preloads {
		db = db.Preload(p)
	}
	return db
}

func (r *BaseRepository[T]) Persist(ctx context.Context, model *T) error {
	return gorm.G[T](r.Db).Create(ctx, model)
}

func (r *BaseRepository[T]) FindOneByID(ctx context.Context, id int64, preloads []string) (*T, error) {
	db := r.preload(r.Db.WithContext(ctx), preloads)
	res, err := gorm.G[T](db).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *BaseRepository[T]) FindOneBy(ctx context.Context, filters []QueryFilter, preloads ...string) (*T, error) {
	db := r.preload(r.Db, preloads)
	for _, f := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", f.Key), f.Value)
	}

	res, err := gorm.G[T](db).First(ctx)
	if err != nil {
		return nil, err
	}

	return &res, err
}

func (r *BaseRepository[T]) FindBy(ctx context.Context, filters []QueryFilter, preloads ...string) ([]T, error) {
	db := r.preload(r.Db, preloads)
	for _, f := range filters {
		db = db.Where(fmt.Sprintf("%s = ?", f.Key), f.Value)
	}

	res, err := gorm.G[T](db).Find(ctx)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (r *BaseRepository[T]) Update(ctx context.Context, id int64, fields map[string]interface{}) error {
	_, err := gorm.G[map[string]interface{}](r.Db).Table(r.Table).Where("id = ?", id).Updates(ctx, fields)
	return err
}

func (r *BaseRepository[T]) Delete(ctx context.Context, id int64) error {
	_, err := gorm.G[T](r.Db).Where("id = ?", id).Delete(ctx)
	return err
}
