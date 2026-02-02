package services

import (
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/repositories"
	"github.com/rangodisco/yhar/internal/api/types/subsonic"
	"github.com/rangodisco/yhar/internal/metadata/types/scrobble"
)

type TrackService struct {
	tRepo *repositories.TrackRepository
}

func NewTrackService(tRepo *repositories.TrackRepository) *TrackService {
	return &TrackService{tRepo: tRepo}
}

// GetTrackByScrobbleInfo tries to find an existing track from its database, based on a subsonic getNowPlaying's entry
func (s *TrackService) GetTrackByScrobbleInfo(entry *subsonic.Entry) (*models.Track, error) {
	track, err := s.tRepo.FindActiveTrackByTitle(entry.Title)
	if err != nil {
		return nil, err
	}
	return track, err
}

// CreateTrackFromMetadata creates a new Track from a scrobble
func (s *TrackService) CreateTrackFromMetadata(info *scrobble.TrackInfo, musicBrainzID string, artists []models.Artist, album models.Album) (*models.Track, error) {
	track := s.buildTrackModel(info, musicBrainzID, artists, album)
	err := s.tRepo.PersistTrack(track)
	if err != nil {
		return nil, err
	}

	return track, nil
}

// buildTrackModel builds a new models.Track based on a scrobble
func (s *TrackService) buildTrackModel(info *scrobble.TrackInfo, musicBrainzID string, artists []models.Artist, album models.Album) *models.Track {
	return &models.Track{
		Title:         info.Title,
		MusicBrainzID: musicBrainzID,
		Artists:       artists,
		Album:         album,
	}
}
