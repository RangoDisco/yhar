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
		trackArtists = append(trackArtists, scrobble.ArtistInfo{
			Name:     artist.Name,
			ImageUrl: "",
		})
	}

	for _, a := range *albums {
		var albumArtists []scrobble.ArtistInfo
		for _, aa := range a.Artists {
			albumArtists = append(albumArtists, scrobble.ArtistInfo{
				Name:     aa.Name,
				ImageUrl: "",
			})
		}
		trackAlbums = append(trackAlbums, scrobble.AlbumInfo{
			Title:   a.Name,
			Artists: albumArtists,
		})
	}

	info.Track = scrobble.TrackInfo{
		Title:   track.Name,
		Artists: trackArtists,
		Album:   trackAlbums,
	}

	return &info, nil
}
