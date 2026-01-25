package services

import (
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/repositories"
	"github.com/rangodisco/yhar/internal/api/types/subsonic"
	"github.com/rangodisco/yhar/internal/metadata/types/scrobble"
)

// GetTrackByScrobbleInfo tries to find an existing track from its database, based on a subsonic getNowPlaying's entry
func GetTrackByScrobbleInfo(entry *subsonic.Entry) (*models.Track, error) {
	track, err := repositories.FindActiveTrackByTitle(entry.Title)
	if err != nil {
		return nil, err
	}
	return track, err
}

// CreateTrackFromMetadata creates a new Track from a scrobble
func CreateTrackFromMetadata(info *scrobble.TrackInfo, musicBrainzID string, artists []models.Artist, album models.Album) (*models.Track, error) {
	track := buildTrackModel(info, musicBrainzID, artists, album)
	err := repositories.PersistTrack(track)
	if err != nil {
		return nil, err
	}

	return track, nil
}

// buildTrackModel builds a new models.Track based on a scrobble
func buildTrackModel(info *scrobble.TrackInfo, musicBrainzID string, artists []models.Artist, album models.Album) *models.Track {
	return &models.Track{
		Title:         info.Title,
		MusicBrainzID: musicBrainzID,
		Artists:       artists,
		Album:         album,
	}
}
