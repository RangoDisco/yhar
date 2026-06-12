package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rangodisco/yhar/internal/api/models"
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
func (s *MetadataService) GetInfoByScrobble(ctx context.Context, MBID, title, artist, album, duration string) (*providers.InfoResponse, error) {
	if MBID == "" && title == "" {
		return nil, fmt.Errorf("%w: at least one MBID or title is required", ErrInvalidScrobble)
	}

	data := &providers.ScrobbleData{
		Title:    title,
		MBID:     MBID,
		Album:    album,
		Artist:   artist,
		Duration: duration,
	}

	info, err := s.enrichMetadata(ctx, data)
	if err != nil {
		// Still create manual track infos for tracks without metadata (unreleased, custom songs, etc.)
		manualInfos, manualErr := s.buildManualProviderInfos(data)
		if manualErr != nil {
			return nil, err
		}
		return manualInfos, nil
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

	// TODO: If MusicBrainz wasn't enough, try all others until data was found
	//if trackInfo == nil {
	//	// TODO:
	//	for _, p := range s.providers {
	//		// Skip MusicBrainz as it was already used
	//		if p.Name() == "musicbrainz" {
	//			continue
	//		}
	//	}
	//}

	// If track info is still nil after all providers were called, return the errors
	if trackInfo == nil {
		return nil, fmt.Errorf("%w, %v", ErrNoMetadataFound, errors.Join(errs...))
	}

	if trackInfo.Album.ImageURL == "" {
		// Fallback to deezer for album image
		dzProvider := s.findProviderByName("deezer")
		img, err := dzProvider.GetAlbumImage(ctx, trackInfo.Album.Title, trackInfo.Album.Artists[0].Name)
		if err == nil {
			trackInfo.Album.ImageURL = img
		}
	}

	// TODO: Maybe handle errs and log ?
	s.addPicturesToArtists(ctx, trackInfo)

	return &providers.InfoResponse{
		Track: *trackInfo,
	}, nil
}

func (s *MetadataService) buildManualProviderInfos(infos *providers.ScrobbleData) (*providers.InfoResponse, error) {
	duration, err := time.ParseDuration(fmt.Sprintf("%ss", infos.Duration))
	if err != nil {
		return nil, err
	}
	artists := []providers.ArtistMetadata{{Name: infos.Artist, SortName: infos.Artist, ImageUrl: "", Genres: make([]string, 0), MBID: ""}}
	return &providers.InfoResponse{
		Track: providers.TrackMetadata{
			Title:   infos.Title,
			Artists: artists,
			Album: providers.AlbumMetadata{
				Title:     infos.Album,
				MBID:      "",
				Artists:   artists,
				AlbumType: string(models.ALBUM),
			},
			Duration: duration,
			ISRC:     "",
			MBID:     "",
		},
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

func (s *MetadataService) addPicturesToArtists(ctx context.Context, trackInfo *providers.TrackMetadata) {
	dProvider := s.findProviderByName("deezer")
	if dProvider == nil {
		return
	}

	for i, trackArtist := range trackInfo.Artists {
		if trackArtist.ImageUrl != "" {
			continue
		}

		img, err := dProvider.GetArtistImage(ctx, "", trackArtist.Name)
		if err != nil {
			continue
		}

		trackInfo.Artists[i].ImageUrl = img
	}

	for i, albumArtists := range trackInfo.Album.Artists {
		alreadyHasImg := false
		// No need to fetch the image if the artists also appears on the track, as it will be found when trying to GetOrCreate
		for _, tAr := range trackInfo.Artists {
			if albumArtists.MBID == tAr.MBID {
				alreadyHasImg = true
				continue
			}
		}

		// Fetch artist image URL
		if !alreadyHasImg {
			img, err := dProvider.GetArtistImage(ctx, "", albumArtists.Name)
			if err != nil {
				continue
			}

			trackInfo.Album.Artists[i].ImageUrl = img
		}
	}
}
