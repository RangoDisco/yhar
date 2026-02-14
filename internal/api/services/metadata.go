package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/rangodisco/yhar/internal/api/providers"
)

type MetadataService struct {
	providers []providers.MetadataProvider
}

var (
	ErrNoMetadataFound = errors.New("no metadata found from any provider")
	ErrInvalidScrobble = errors.New("invalid scrobble data")
)

func NewMetadataService(
	providers []providers.MetadataProvider,
) *MetadataService {
	return &MetadataService{providers: providers}
}

// GetInfoByScrobble fetches metadata from multiple providers and formats it into a standardized providers.InfoResponse
func (s *MetadataService) GetInfoByScrobble(ctx context.Context, MBID, title string) (*providers.InfoResponse, error) {
	if MBID == "" && title == "" {
		return nil, fmt.Errorf("%w: at least one MBID or title is required", ErrInvalidScrobble)
	}

	info, err := s.enrichMetadata(ctx, &providers.ScrobbleData{
		Title: title,
		MBID:  MBID,
	})

	if err != nil {
		return nil, fmt.Errorf("unable to enrich metadata: %w", err)
	}

	return info, nil
}

func (s *MetadataService) enrichMetadata(ctx context.Context, infos *providers.ScrobbleData) (*providers.InfoResponse, error) {
	var trackInfo *providers.TrackMetadata
	var errs []error

	// First always try MusicBrainz as it's the most complete
	mBProvider := s.findProviderByName("musicbrainz")
	if mBProvider != nil {
		track, err := mBProvider.GetTrackByInfos(ctx, *infos)
		if err == nil && track != nil {
			trackInfo = track
		} else {
			errs = append(errs, fmt.Errorf("muscibrainz provider: %w", err))
		}
	}

	// If MusicBrainz wasn't enough, try all others until data was found
	if trackInfo == nil {
		// TODO:
		for _, p := range s.providers {
			// Skip MusicBrainz as it was already used
			if p.Name() == "muscibrainz" {
				continue
			}
		}
	}

	// If track info is still nil after all providers were called, return the errors
	if trackInfo == nil {
		return nil, fmt.Errorf("%w, %v", ErrNoMetadataFound, errors.Join(errs...))
	}

	// Maybe handle errs and log ?
	s.addPicturesToArtists(ctx, trackInfo)

	return &providers.InfoResponse{
		Track: *trackInfo,
	}, nil

}

func (s *MetadataService) findProviderByName(name string) providers.MetadataProvider {
	for _, p := range s.providers {
		if p.Name() == name {
			return p
		}
	}
	return nil
}

// TODO: handle range over pointer ?
func (s *MetadataService) addPicturesToArtists(ctx context.Context, trackInfo *providers.TrackMetadata) {
	for _, trackArtists := range trackInfo.Artists {
		if trackArtists.ImageUrl != "" {
			continue
		}

		// FETCH artist image URL
	}

	for _, albumArtists := range trackInfo.Album.Artists {
		// No need to fetch the image if the artists also appears on the track, as it will be found when trying to GetOrCreate
		for _, tAr := range trackInfo.Artists {
			if albumArtists.MBID == tAr.MBID {
				continue
			}
			// Fetch artist image URL
		}
	}
}
