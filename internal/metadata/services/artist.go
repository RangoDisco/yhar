package services

import (
	"github.com/rangodisco/yhar/internal/metadata/models"
	"github.com/rangodisco/yhar/internal/metadata/types/scrobble"
)

type ArtistService struct {
}

func NewArtistService() *ArtistService {
	return &ArtistService{}
}

func (s *ArtistService) FormatArtistToScrobbleInfo(artist *models.Artist) *scrobble.ArtistInfo {
	imgUrl := ""
	for _, image := range artist.Images {
		if image.Url != "" {
			imgUrl = image.Url
			break
		}
	}

	var genres []string
	for _, genre := range artist.Genres {
		genres = append(genres, genre.Name)
	}

	return &scrobble.ArtistInfo{
		Name:     artist.Name,
		ImageUrl: imgUrl,
		Genres:   genres,
	}
}
