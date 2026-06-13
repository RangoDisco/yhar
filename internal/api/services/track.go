package services

import (
	"context"

	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/providers"
	"github.com/rangodisco/yhar/internal/api/repositories"
)

type TrackService struct {
	repo *repositories.TrackRepository
}

func NewTrackService(repo *repositories.TrackRepository) *TrackService {
	return &TrackService{repo: repo}
}

// GetByScrobbleInfo tries to find an existing models.Track from the database, based on a subsonic.Entry
func (s *TrackService) GetByScrobbleInfo(ctx context.Context, entry *UnifiedScrobbleEntry) (*models.Track, error) {
	filter := repositories.QueryFilter{Key: "title", Value: entry.Title}
	if entry.MusicBrainzID != nil {
		filter = repositories.QueryFilter{Key: "music_brainz_id", Value: entry.MusicBrainzID}
	}

	track, err := s.repo.FindOneBy(ctx, []repositories.QueryFilter{filter}, "Artists.Picture", "Album.Picture")
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
		AlbumID:       album.ID,
		Duration:      info.Duration,
	}

	err := s.repo.Persist(ctx, track)
	if err != nil {
		return nil, err
	}

	return track, nil
}
