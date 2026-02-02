package services

import (
	"github.com/rangodisco/yhar/internal/metadata/models"
	"github.com/rangodisco/yhar/internal/metadata/types/scrobble"
)

type AlbumService struct {
	arService *ArtistService
}

func NewAlbumService(arService *ArtistService) *AlbumService {
	return &AlbumService{arService: arService}
}

func (s *AlbumService) FormatAlbumToScrobbleInfo(album *models.Album) *scrobble.AlbumInfo {
	var albumArtists []scrobble.ArtistInfo
	for _, aa := range album.Artists {
		artistInfo := s.arService.FormatArtistToScrobbleInfo(&aa)
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
		Title:     album.Name,
		ImageUrl:  imgUrl,
		Artists:   albumArtists,
		AlbumType: album.AlbumType,
	}
}
