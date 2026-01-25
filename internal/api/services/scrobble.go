package services

import (
	"errors"

	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/repositories"
	"github.com/rangodisco/yhar/internal/api/types/subsonic"
	"github.com/rangodisco/yhar/internal/metadata/services"
	"github.com/rangodisco/yhar/internal/metadata/types/scrobble"
)

func HandleNewScrobble(entry subsonic.Entry) (*models.Scrobble, error) {
	user, err := GetOrCreateUser(entry.Username)
	if err != nil {
		return nil, err
	}

	var track *models.Track

	// See if track already exist in database
	track, _ = GetTrackByScrobbleInfo(&entry)

	if track == nil {
		track, err = getOrCreateScrobbleContents(&entry)
		if err != nil {
			return nil, err
		}
	}

	// Create and persist new scrobble
	model := buildScrobbleModel(track.ID, user.ID)

	err = repositories.PersistScrobble(model)
	if err != nil {
		return nil, err
	}
	return model, nil
}

// getOrCreateScrobbleContents gets or create all content (track, album, artists) related to the scrobble
func getOrCreateScrobbleContents(entry *subsonic.Entry) (*models.Track, error) {
	// Otherwise fetch metadata and create it
	data, err := GetTrackMetadata(entry)
	if err != nil {
		return nil, err
	}

	artists := ProcessScrobbleArtists(data.Track.Artists)
	if len(artists) == 0 {
		return nil, errors.New("no artists found")
	}

	// Get album
	album := ProcessScrobbleAlbums(data.Track.Albums[0])

	// Create track with everything
	track, err := CreateTrackFromMetadata(&data.Track, entry.MusicBrainzID, artists, album)
	if err != nil {
		return nil, err
	}

	return track, nil
}

// ProcessScrobbleArtists takes all artists given, and find or create their picture and themselves
func ProcessScrobbleArtists(sArtists []scrobble.ArtistInfo) []models.Artist {
	// Create all artists related to the track
	var artists []models.Artist
	for _, artistInfo := range sArtists {
		artist, err := GetOrCreateArtist(artistInfo)
		if err != nil {
			continue
		}
		artists = append(artists, *artist)
	}

	return artists
}

// ProcessScrobbleAlbums takes all albums given, and find or create their artists, picture and themselves
func ProcessScrobbleAlbums(sAlbum scrobble.AlbumInfo) models.Album {
	artists := ProcessScrobbleArtists(sAlbum.Artists)
	album, _ := GetOrCreateAlbum(sAlbum, artists)

	return *album
}

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

func buildScrobbleModel(trackID, userID int64) *models.Scrobble {
	return &models.Scrobble{
		Origin:  models.SUBSONIC,
		TrackID: trackID,
		UserID:  userID,
	}
}
