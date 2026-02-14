package repositories

import (
	"context"

	"github.com/rangodisco/yhar/internal/api/models"
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

func (r *ArtistRepository) FindActiveArtistByName(ctx context.Context, name string) (*models.Artist, error) {
	var a models.Artist

	err := r.Db.WithContext(ctx).Where("name = ?", name).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, err
}

func (r *ArtistRepository) PersistArtist(ctx context.Context, a *models.Artist) error {
	res := r.Db.WithContext(ctx).Create(&a)
	return res.Error
}
