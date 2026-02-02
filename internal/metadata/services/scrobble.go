package services

import (
	"github.com/rangodisco/yhar/internal/metadata/models"
	"github.com/rangodisco/yhar/internal/metadata/repositories"
	"github.com/rangodisco/yhar/internal/metadata/types/scrobble"
)

type ScrobbleService struct {
	tRepo     *repositories.TrackRepository
	alRepo    *repositories.AlbumRepository
	arService *ArtistService
	alService *AlbumService
}

func NewScrobbleService(
	tRepo *repositories.TrackRepository,
	alRepo *repositories.AlbumRepository,
	arService *ArtistService,
	alService *AlbumService,
) *ScrobbleService {
	return &ScrobbleService{tRepo: tRepo, alRepo: alRepo, arService: arService, alService: alService}
}

func (s *ScrobbleService) GetInfoByScrobble(scrobble scrobble.InfoRequest) (*scrobble.InfoResponse, error) {
	t, err := s.tRepo.FindTrackInfoByScrobble(scrobble)
	if err != nil {
		return nil, err
	}

	a, err := s.alRepo.FindAlbumById(t.AlbumID)
	if err != nil {
		return nil, err
	}

	info, err := s.formatToInfoByScrobbleResponse(t, a)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (s *ScrobbleService) formatToInfoByScrobbleResponse(track *models.Track, albums *[]models.Album) (*scrobble.InfoResponse, error) {
	var info scrobble.InfoResponse
	var trackArtists []scrobble.ArtistInfo
	var trackAlbums []scrobble.AlbumInfo

	for _, artist := range track.Artists {
		arInfo := s.arService.FormatArtistToScrobbleInfo(&artist)
		trackArtists = append(trackArtists, *arInfo)
	}

	for _, a := range *albums {
		alInfo := s.alService.FormatAlbumToScrobbleInfo(&a)
		trackAlbums = append(trackAlbums, *alInfo)
	}

	info.Track = scrobble.TrackInfo{
		Title:   track.Name,
		Artists: trackArtists,
		Albums:  trackAlbums,
	}

	return &info, nil
}
