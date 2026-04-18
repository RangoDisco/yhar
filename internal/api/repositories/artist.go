package repositories

import (
	"context"
	"fmt"

	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/types/filters"
	"gorm.io/gorm"
)

type ArtistRepository struct {
	Db *gorm.DB
}

func NewArtistRepository(Db *gorm.DB) *ArtistRepository {
	return &ArtistRepository{
		Db: Db,
	}
}

func (r *ArtistRepository) FindActiveByFilters(ctx context.Context, filters []filters.QueryFilter) (*models.Artist, error) {
	var a models.Artist

	if len(filters) == 0 {
		return nil, fmt.Errorf("filters are empty")
	}

	query := r.Db.WithContext(ctx)

	for _, f := range filters {
		query = query.Where(fmt.Sprintf("%s = ?", f.Key), f.Value)
	}

	err := query.Preload("Picture").First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, err
}

func (r *ArtistRepository) Persist(ctx context.Context, a *models.Artist) error {
	res := r.Db.WithContext(ctx).Create(&a)
	return res.Error
}

func (r *ArtistRepository) Update(ctx context.Context, a *models.Artist, updates map[string]interface{}) error {
	_, err := gorm.G[map[string]interface{}](r.Db).Table("artists").Where("id = ?", a.ID).Updates(ctx, updates)
	return err
}
