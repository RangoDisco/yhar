package services

import (
	"context"

	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/providers"
	"github.com/rangodisco/yhar/internal/api/repositories"
	filters "github.com/rangodisco/yhar/internal/api/types/filters"
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
	var queryFilters []filters.QueryFilter

	if info.MBID != "" {
		queryFilters = append(queryFilters, filters.QueryFilter{Key: "music_brainz_id", Value: info.MBID})
	} else {
		queryFilters = append(queryFilters, filters.QueryFilter{Key: "name", Value: info.Name})
	}

	existingArtist, err := s.repo.FindActiveByFilters(ctx, queryFilters)
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

	err = s.repo.Persist(ctx, model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

// scrobbleInfoToArtistModel builds a new models.Artist based on a scrobble
func scrobbleInfoToArtistModel(info providers.ArtistMetadata, img *models.Image) *models.Artist {
	return &models.Artist{
		Name:          info.Name,
		PictureID:     img.ID,
		MusicBrainzID: info.MBID,
	}
}
