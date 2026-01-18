package services

import (
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/repositories"
	"github.com/rangodisco/yhar/internal/metadata/types/scrobble"
)

func GetOrCreateAlbum(info scrobble.AlbumInfo, artists []models.Artist) (*models.Album, error) {
	existingAlbum, err := repositories.FindActiveAlbumByTitle(info.Title)
	if err == nil {
		return existingAlbum, nil
	}
	img, _ := GetOrCreateImage(info.ImageUrl)
	model := scrobbleInfoToAlbumModel(info, artists, img)

	newAlbum, err := repositories.CreateAlbum(model)
	if err != nil {
		return nil, err
	}
	return newAlbum, nil
}

func scrobbleInfoToAlbumModel(info scrobble.AlbumInfo, artists []models.Artist, img *models.Image) *models.Album {
	return &models.Album{
		Title:     info.Title,
		Artists:   artists,
		PictureID: img.ID,
	}
}
