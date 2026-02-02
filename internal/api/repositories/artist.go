package repositories

import (
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

func (r *ArtistRepository) FindActiveArtistByName(name string) (*models.Artist, error) {
	var a models.Artist

	err := r.Db.Where("name = ?", name).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, err
}

func (r *ArtistRepository) PersistArtist(a *models.Artist) error {
	res := r.Db.Create(&a)
	return res.Error
}
