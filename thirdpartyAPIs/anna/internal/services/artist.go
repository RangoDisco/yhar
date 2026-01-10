package services

import (
	"github.com/rangodisco/yhar/pkg/types/anna/scrobble"
	"github.com/rangodisco/yhar/thirdpartyAPIs/anna/internal/models"
)

func FormatArtistToScrobbleInfo(artist *models.Artist) *scrobble.ArtistInfo {
	imgUrl := ""
	for _, image := range artist.Images {
		if image.Url != "" {
			imgUrl = image.Url
			break
		}
	}

	return &scrobble.ArtistInfo{
		Name:     artist.Name,
		ImageUrl: imgUrl,
	}
}
