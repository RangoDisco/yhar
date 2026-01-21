package services

import (
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/repositories"
	"github.com/rangodisco/yhar/internal/api/types/subsonic"
	"github.com/rangodisco/yhar/internal/metadata/types/scrobble"
)

func GetTrackByScrobbleInfo(entry *subsonic.Entry) (*models.Track, error) {
	track, err := repositories.FindActiveTrackByTitle(entry.Title)
	if err != nil {
		return nil, err
	}
	return track, err
}

func CreateTrackFromMetadata(info *scrobble.TrackInfo, artists []models.Artist, album models.Album) (*models.Track, error) {
	track := buildTrackModel(info, artists, album)
	err := repositories.PersistTrack(track)
	if err != nil {
		return nil, err
	}

	return track, nil
}

func buildTrackModel(info *scrobble.TrackInfo, artists []models.Artist, album models.Album) *models.Track {
	return &models.Track{
		Title:   info.Title,
		Artists: artists,
		Album:   album,
	}
}
