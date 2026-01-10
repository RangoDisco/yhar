package services

import (
	"github.com/rangodisco/yhar/pkg/types/anna/scrobble"
	"github.com/rangodisco/yhar/thirdpartyAPIs/anna/internal/models"
)

func FormatAlbumToScrobbleInfo(album *models.Album) *scrobble.AlbumInfo {
	var albumArtists []scrobble.ArtistInfo
	for _, aa := range album.Artists {
		artistInfo := FormatArtistToScrobbleInfo(&aa)
		albumArtists = append(albumArtists, *artistInfo)
	}

	imgUrl := ""

	for _, image := range album.Images {
		if image.Url != "" {
			imgUrl = image.Url
			break
		}
	}

	return &scrobble.AlbumInfo{
		Title:    album.Name,
		ImageUrl: imgUrl,
		Artists:  albumArtists,
	}
}
