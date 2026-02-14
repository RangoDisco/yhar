package services

import (
	"context"

	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/providers"
	"github.com/rangodisco/yhar/internal/api/repositories"
)

type ArtistService struct {
	repo  *repositories.ArtistRepository
	image *ImageService
	genre *GenreService
}

func NewArtistService(repo *repositories.ArtistRepository, image *ImageService, genre *GenreService) *ArtistService {
	return &ArtistService{repo: repo, image: image, genre: genre}
}

// GetOrCreate tries to fetch or create a models.Artist if it doesn't exist
func (s *ArtistService) GetOrCreate(ctx context.Context, info providers.ArtistMetadata) (*models.Artist, error) {
	// TODO: find by MBID and not name
	existingArtist, err := s.repo.FindActiveArtistByName(ctx, info.Name)
	if err == nil && existingArtist.Name != "" {
		return existingArtist, err
	}

	img, _ := s.image.GetOrCreate(ctx, info.ImageUrl)

	//// Add all genres needed for the future model
	//var genres []models.Genre
	//for _, genreInfo := range info.Genres {
	//	genre, err := s.genre.GetOrCreateGenre(ctx, genreInfo)
	//	if err != nil {
	//		// We don't want to stop the whole request just for a missing genre
	//		continue
	//	}
	//	genres = append(genres, *genre)
	//}

	// Build the model object from all the infos
	model := scrobbleInfoToArtistModel(info, img)

	err = s.repo.PersistArtist(ctx, model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

// scrobbleInfoToArtistModel builds a new models.Artist based on a scrobble
func scrobbleInfoToArtistModel(info providers.ArtistMetadata, img *models.Image) *models.Artist {
	return &models.Artist{
		Name:      info.Name,
		PictureID: img.ID,
	}
}
