package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/rangodisco/yhar/internal/api/dto"
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/providers"
	"github.com/rangodisco/yhar/internal/api/repositories"
	filters "github.com/rangodisco/yhar/internal/api/types/filters"
	"gorm.io/gorm"
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

	existingArtist, err := s.Get(ctx, queryFilters)
	if err == nil && existingArtist != nil {
		return existingArtist, err
	}

	var img *models.Image
	if info.ImageUrl != "" {
		img, err = s.image.GetOrCreate(ctx, info.ImageUrl)
		if err != nil {
			return nil, fmt.Errorf("unable to get or create img: %w", err)
		}
	}

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
	a := scrobbleInfoToArtistModel(info, img)

	err = s.repo.Persist(ctx, a)
	if err != nil {
		// could happen if another routine inserted the same artist first
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return s.repo.FindActiveByFilters(ctx, queryFilters)
		}
		return nil, err
	}

	return a, nil
}

// Get finds a models.Artist by filters
func (s *ArtistService) Get(ctx context.Context, filters []filters.QueryFilter) (*models.Artist, error) {
	return s.repo.FindActiveByFilters(ctx, filters)
}

// Update partially updates a models.Artist with given dto.UpdateArtistInput
func (s *ArtistService) Update(ctx context.Context, a *models.Artist, input dto.UpdateArtistInput) (*models.Artist, error) {
	updates := make(map[string]interface{})

	if input.Name != nil {
		updates["name"] = *input.Name
	}

	if input.ImageID != nil {
		updates["picture_id"] = *input.ImageID
	}

	err := s.repo.Update(ctx, a, updates)
	if err != nil {
		return nil, fmt.Errorf("unable to update artist: %w", err)
	}

	return a, nil
}

// scrobbleInfoToArtistModel builds a new models.Artist based on a scrobble
func scrobbleInfoToArtistModel(info providers.ArtistMetadata, img *models.Image) *models.Artist {
	model := models.Artist{
		Name:          info.Name,
		MusicBrainzID: info.MBID,
	}

	if img != nil {
		model.PictureID = &img.ID
	}

	return &model
}
