package repositories

import (
	"github.com/rangodisco/yhar/pkg/types/anna/scrobble"
	"github.com/rangodisco/yhar/thirdpartyAPIs/anna/config/database"
	"github.com/rangodisco/yhar/thirdpartyAPIs/anna/internal/models"
)

func FindTrackInfoByScrobble(scrobble scrobble.InfoRequest) (*models.Track, error) {
	var t models.Track
	err := database.GetDB().Preload("Artists").Where("name = ?", scrobble.Title).Order("popularity").Last(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}
