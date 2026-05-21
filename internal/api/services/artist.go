package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/rangodisco/yhar/internal/api/dto"
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/providers"
	"github.com/rangodisco/yhar/internal/api/repositories"
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
	filter := repositories.QueryFilter{Key: "name", Value: info.Name}

	if info.MBID != "" {
		filter = repositories.QueryFilter{
			Key: "music_brainz_id", Value: info.MBID,
		}
	}

	existingArtist, err := s.repo.FindOneBy(ctx, []repositories.QueryFilter{filter}, "Picture")
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
			return s.repo.FindOneBy(ctx, []repositories.QueryFilter{filter}, "Picture")
		}
		return nil, err
	}

	return a, nil
}

// GetByID finds a models.Artist by ID
func (s *ArtistService) GetByID(ctx context.Context, id int64) (*models.Artist, error) {
	return s.repo.FindOneByID(ctx, id, "Picture")
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

	err := s.repo.Update(ctx, a.ID, updates)
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
