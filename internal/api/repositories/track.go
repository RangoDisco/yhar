package repositories

import (
	"github.com/rangodisco/yhar/internal/api/config/database"
	"github.com/rangodisco/yhar/internal/api/models"
)

func FindActiveTrackByTitle(title string) (*models.Track, error) {
	var t models.Track

	// TODO: handle multiple track with same name (check for albums/artists)
	err := database.GetDB().Preload("Artists.Picture").Preload("Album.Picture").Where("title = ?", title).First(&t).Error
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func PersistTrack(track *models.Track) error {
	res := database.GetDB().Create(&track)

	return res.Error
}
