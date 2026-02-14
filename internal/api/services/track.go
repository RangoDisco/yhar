package services

import (
	"context"

	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/providers"
	"github.com/rangodisco/yhar/internal/api/repositories"
	"github.com/rangodisco/yhar/internal/api/types/subsonic"
)

type TrackService struct {
	repo *repositories.TrackRepository
}

func NewTrackService(repo *repositories.TrackRepository) *TrackService {
	return &TrackService{repo: repo}
}

// GetByScrobbleInfo tries to find an existing models.Track from the database, based on a subsonic.Entry
func (s *TrackService) GetByScrobbleInfo(ctx context.Context, entry *subsonic.Entry) (*models.Track, error) {
	// TODO: handle if mbid is present
	track, err := s.repo.FindActiveByTitle(ctx, entry.Title)
	if err != nil {
		return nil, err
	}
	return track, err
}

// CreateFromMetadata creates a new models.Track from a providers.TrackMetadata
func (s *TrackService) CreateFromMetadata(ctx context.Context, info *providers.TrackMetadata, artists []models.Artist, album models.Album) (*models.Track, error) {
	track := &models.Track{
		Title:         info.Title,
		MusicBrainzID: info.MBID,
		Artists:       artists,
		Album:         album,
	}

	err := s.repo.PersistTrack(ctx, track)
	if err != nil {
		return nil, err
	}

	return track, nil
}
