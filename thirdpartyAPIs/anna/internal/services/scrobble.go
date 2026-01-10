package services

import (
	"github.com/rangodisco/yhar/pkg/types/anna/scrobble"
	"github.com/rangodisco/yhar/thirdpartyAPIs/anna/internal/models"
	"github.com/rangodisco/yhar/thirdpartyAPIs/anna/internal/repositories"
)

func GetInfoByScrobble(scrobble scrobble.InfoRequest) (*scrobble.InfoResponse, error) {
	t, err := repositories.FindTrackInfoByScrobble(scrobble)
	if err != nil {
		return nil, err
	}

	a, err := repositories.FindAlbumById(t.AlbumID)
	if err != nil {
		return nil, err
	}

	info, err := formatToInfoByScrobbleResponse(t, a)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func formatToInfoByScrobbleResponse(track *models.Track, albums *[]models.Album) (*scrobble.InfoResponse, error) {
	var info scrobble.InfoResponse
	var trackArtists []scrobble.ArtistInfo
	var trackAlbums []scrobble.AlbumInfo

	for _, artist := range track.Artists {
		arInfo := FormatArtistToScrobbleInfo(&artist)
		trackArtists = append(trackArtists, *arInfo)
	}

	for _, a := range *albums {
		alInfo := FormatAlbumToScrobbleInfo(&a)
		trackAlbums = append(trackAlbums, *alInfo)
	}

	info.Track = scrobble.TrackInfo{
		Title:   track.Name,
		Artists: trackArtists,
		Album:   trackAlbums,
	}

	return &info, nil
}
