package services

import (
	"context"

	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/providers"
	"github.com/rangodisco/yhar/internal/api/repositories"
	"github.com/rangodisco/yhar/internal/api/types/filters"
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
	var queryFilters []filters.QueryFilter
	if entry.MusicBrainzID != "" {
		queryFilters = append(queryFilters, filters.QueryFilter{Key: "music_brainz_id", Value: entry.MusicBrainzID})
	} else {
		queryFilters = append(queryFilters, filters.QueryFilter{Key: "title", Value: entry.Title})
	}

	track, err := s.repo.FindActiveByFilter(ctx, queryFilters)
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
