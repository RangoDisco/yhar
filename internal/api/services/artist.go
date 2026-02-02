package services

import (
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/repositories"
	"github.com/rangodisco/yhar/internal/metadata/types/scrobble"
)

type ArtistService struct {
	aRepo    *repositories.ArtistRepository
	iService *ImageService
	gService *GenreService
}

func NewArtistService(aRepo *repositories.ArtistRepository, iService *ImageService, gService *GenreService) *ArtistService {
	return &ArtistService{aRepo: aRepo, iService: iService, gService: gService}
}

// GetOrCreateArtist tries to fetch or create an artist if it doesn't exist
func (s *ArtistService) GetOrCreateArtist(info scrobble.ArtistInfo) (*models.Artist, error) {
	existingArtist, err := s.aRepo.FindActiveArtistByName(info.Name)
	if err == nil && existingArtist.Name != "" {
		return existingArtist, err
	}

	img, _ := s.iService.GetOrCreateImage(info.ImageUrl)

	// Add all genres needed for the future model
	var genres []models.Genre
	for _, genreInfo := range info.Genres {
		genre, err := s.gService.GetOrCreateGenre(genreInfo)
		if err != nil {
			// We don't want to stop the whole request just for a missing genre
			continue
		}
		genres = append(genres, *genre)
	}

	// Build the model object from all the infos
	model := scrobbleInfoToArtistModel(info, img, genres)

	err = s.aRepo.PersistArtist(model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

// scrobbleInfoToArtistModel builds a new models.Artist based on a scrobble
func scrobbleInfoToArtistModel(info scrobble.ArtistInfo, img *models.Image, genres []models.Genre) *models.Artist {
	return &models.Artist{
		Name:      info.Name,
		PictureID: img.ID,
		Genres:    genres,
	}
}
