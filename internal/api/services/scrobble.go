package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/providers"
	"github.com/rangodisco/yhar/internal/api/repositories"
	"gorm.io/gorm"
)

type ScrobbleService struct {
	repo     *repositories.ScrobbleRepository
	user     *UserService
	track    *TrackService
	artist   *ArtistService
	album    *AlbumService
	metadata *MetadataService
}

func NewScrobbleService(
	repo *repositories.ScrobbleRepository,
	user *UserService,
	track *TrackService,
	artist *ArtistService,
	album *AlbumService,
	metadata *MetadataService,
) *ScrobbleService {
	return &ScrobbleService{
		repo:     repo,
		user:     user,
		track:    track,
		artist:   artist,
		album:    album,
		metadata: metadata,
	}
}

type UnifiedScrobbleEntry struct {
	Username      string    `json:"username"`
	Title         string    `json:"title"`
	Artist        string    `json:"artist"`
	Album         string    `json:"album"`
	MusicBrainzID string    `json:"musicbrainz_id"`
	ListenedAt    time.Time `json:"listened_at"`
}

// HandleNewScrobble takes a subsonic getNowPlaying entry, fetches/creates associated content and persists a new Scrobble
func (s *ScrobbleService) HandleNewScrobble(ctx context.Context, entry *UnifiedScrobbleEntry) (*models.Scrobble, error) {
	// TODO: validate entry
	// err := validateEntry(&entry)

	// Get user linked to the scrobble
	user, err := s.user.GetOrCreateUser(ctx, entry.Username)
	if err != nil {
		return nil, err
	}

	var t *models.Track

	// See if track already exists in database
	t, err = s.getOrCreateTrack(ctx, entry)
	if err != nil {
		return nil, err
	}

	var listenedAt time.Time
	if !entry.ListenedAt.IsZero() {
		listenedAt = entry.ListenedAt
	} else {
		listenedAt = time.Now()
	}

	// Create and persist new scrobble
	scrobble := &models.Scrobble{
		Origin:      models.SUBSONIC,
		Track:       *t,
		TrackID:     t.ID,
		UserID:      user.ID,
		ScrobbledAt: listenedAt,
	}

	// Ensure scrobble wasn't already created
	existingScrobble, err := s.repo.FindByTrackAndTimestamp(ctx, t.ID, user.ID, listenedAt)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("unable to query db: %w", err)
	}

	if existingScrobble != nil {
		return nil, fmt.Errorf("scrobble already exists")
	}

	err = s.repo.PersistScrobble(ctx, scrobble)
	if err != nil {
		return nil, err
	}
	return scrobble, nil
}

// getOrCreateTrack finds or create all content (track, album, artists) related to the scrobble
func (s *ScrobbleService) getOrCreateTrack(ctx context.Context, entry *UnifiedScrobbleEntry) (*models.Track, error) {
	// First, try to find existing track
	existingTrack, err := s.track.GetByScrobbleInfo(ctx, entry)

	// TODO: check if err is not found with custom type
	if err == nil {
		return existingTrack, nil
	}

	metadata, err := s.metadata.GetInfoByScrobble(ctx, entry.MusicBrainzID, entry.Title, entry.Artist, entry.Album)
	if err != nil {
		return nil, err
	}

	// Process and create all associated models
	artists, err := s.processScrobbleArtists(ctx, metadata.Track.Artists)
	if err != nil || len(artists) == 0 {
		return nil, fmt.Errorf("unable to process found artists :%w", err)
	}

	// Get album
	album, err := s.processScrobbleAlbums(ctx, metadata.Track.Album)
	if err != nil {
		return nil, err
	}

	// Create track with everything
	track, err := s.track.CreateFromMetadata(ctx, &metadata.Track, artists, *album)
	if err != nil {
		// could happen if another routine inserted the same track first
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return s.track.GetByScrobbleInfo(ctx, entry)
		}
		return nil, err
	}

	return track, nil
}

// processScrobbleArtists finds or creates all models.Artist associated to a providers.TrackMetadata/providers.AlbumMetadata
func (s *ScrobbleService) processScrobbleArtists(ctx context.Context, sArtists []providers.ArtistMetadata) ([]models.Artist, error) {
	var artists []models.Artist
	var errs []error
	for _, artistInfo := range sArtists {
		artist, err := s.artist.GetOrCreate(ctx, artistInfo)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		artists = append(artists, *artist)
	}

	// TODO: define if needed
	if len(artists) == 0 {
		return nil, fmt.Errorf("unable to retrieve artist %v", errors.Join(errs...))
	}

	return artists, nil
}

// processScrobbleAlbums takes all albums given, and find or create their artists, picture and themselves
func (s *ScrobbleService) processScrobbleAlbums(ctx context.Context, sAlbum providers.AlbumMetadata) (*models.Album, error) {
	artists, err := s.processScrobbleArtists(ctx, sAlbum.Artists)
	if err != nil {
		return nil, err
	}
	album, err := s.album.GetOrCreateAlbum(ctx, sAlbum, artists)
	if err != nil {
		return nil, err
	}

	return album, err
}
