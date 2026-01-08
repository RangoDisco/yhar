package repositories

import (
	"github.com/rangodisco/yhar/pkg/types/anna/track"
	"github.com/rangodisco/yhar/thirdpartyAPIs/anna/config/database"
	"github.com/rangodisco/yhar/thirdpartyAPIs/anna/internal/models"
)

func GetTrackInfoByScrobble(scrobble track.InfoByScrobbleRequest) (*models.Track, error) {
	var t models.Track
	err := database.GetDB().Where("name = ?", scrobble.Title).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}
