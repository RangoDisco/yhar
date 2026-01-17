package services

import (
	"github.com/rangodisco/yhar/internal/api/types/subsonic"
	"github.com/rangodisco/yhar/internal/metadata/services"
	"github.com/rangodisco/yhar/internal/metadata/types/scrobble"
)

// GetTrackMetadata Fetch the current playing track from all setup providers (only subsonic api for now)
// then fetches associated metadata from providers (only local db for now)
func GetTrackMetadata(entry *subsonic.Entry) (*scrobble.InfoResponse, error) {
	scrobbleRequest := &scrobble.InfoRequest{
		Title:  entry.Title,
		Album:  entry.Album,
		Artist: entry.Artist,
	}

	aRes, err := services.GetInfoByScrobble(*scrobbleRequest)
	if err != nil {
		return nil, err
	}

	return aRes, nil
}
