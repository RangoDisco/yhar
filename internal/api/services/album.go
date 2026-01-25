package services

import (
	"fmt"
	"strings"

	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/repositories"
	"github.com/rangodisco/yhar/internal/metadata/types/scrobble"
)

// GetOrCreateAlbum tries to fetch or create an album if it doesn't exist
func GetOrCreateAlbum(info scrobble.AlbumInfo, artists []models.Artist) (*models.Album, error) {
	existingAlbum, err := repositories.FindActiveAlbumByTitle(info.Title)
	if err == nil {
		return existingAlbum, nil
	}

	img, _ := GetOrCreateImage(info.ImageUrl)
	model, err := scrobbleInfoToAlbumModel(info, artists, img)
	if err != nil {
		return nil, err
	}

	err = repositories.PersistAlbum(model)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func parseAlbumType(at string) (*models.AlbumType, error) {
	m := map[models.AlbumType]struct{}{
		models.ALBUM:       {},
		models.EP:          {},
		models.SINGLE:      {},
		models.COMPILATION: {},
	}
	albumType := models.AlbumType(strings.ToUpper(at))

	_, ok := m[albumType]
	if !ok {
		return nil, fmt.Errorf("unable to parse %s as AlbumType", at)
	}

	return &albumType, nil
}

// scrobbleInfoToAlbumModel build a new models.Album based on a scrobble
func scrobbleInfoToAlbumModel(info scrobble.AlbumInfo, artists []models.Artist, img *models.Image) (*models.Album, error) {
	at, err := parseAlbumType(info.AlbumType)

	if err != nil {
		return nil, err
	}

	return &models.Album{
		Title:     info.Title,
		Artists:   artists,
		PictureID: img.ID,
		Type:      *at,
	}, nil
}
