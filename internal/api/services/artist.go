package services

import (
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/repositories"
	"github.com/rangodisco/yhar/internal/metadata/types/scrobble"
)

func GetOrCreateArtist(info scrobble.ArtistInfo, img *models.Image) (*models.Artist, error) {
	existingArtist, err := repositories.FindActiveArtistByName(info.Name)
	if err == nil && existingArtist.Name != "" {
		return existingArtist, err
	}

	model := scrobbleInfoToArtistModel(info, img)

	newArtist, err := repositories.PersistArtist(model)
	if err != nil {
		return nil, err
	}

	return newArtist, nil
}

func scrobbleInfoToArtistModel(info scrobble.ArtistInfo, img *models.Image) *models.Artist {
	return &models.Artist{
		Name:      info.Name,
		PictureID: img.ID,
	}
}
