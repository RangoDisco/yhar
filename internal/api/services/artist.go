package services

import (
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/repositories"
	"github.com/rangodisco/yhar/internal/metadata/types/scrobble"
)

func GetOrCreateArtist(info scrobble.ArtistInfo) (*models.Artist, error) {
	existingArtist, err := repositories.FindActiveArtistByName(info.Name)
	if err == nil && existingArtist.Name != "" {
		return existingArtist, err
	}

	model := createModelFromScrobbleInfo(info)
	newArtist, err := repositories.CreateArtist(model)
	if err != nil {
		return nil, err
	}

	return newArtist, nil
}

func createModelFromScrobbleInfo(info scrobble.ArtistInfo) *models.Artist {
	image, _ := GetOrCreateImage(info.ImageUrl)
	return &models.Artist{
		Name:      info.Name,
		PictureID: image.ID,
	}
}
