package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/rangodisco/yhar/internal/api/dto"
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/providers"
	"github.com/rangodisco/yhar/internal/api/repositories"
	"github.com/rangodisco/yhar/internal/api/types/filters"
	"gorm.io/gorm"
)

type AlbumService struct {
	repo  *repositories.AlbumRepository
	image *ImageService
}

func NewAlbumService(repo *repositories.AlbumRepository, image *ImageService) *AlbumService {
	return &AlbumService{
		repo:  repo,
		image: image,
	}
}

// GetOrCreateAlbum tries to fetch or create an album if it doesn't exist
func (s *AlbumService) GetOrCreateAlbum(ctx context.Context, info providers.AlbumMetadata, artists []models.Artist) (*models.Album, error) {
	var queryFilters []filters.QueryFilter

	if info.MBID != "" {
		queryFilters = append(queryFilters, filters.QueryFilter{
			Key: "music_brainz_id", Value: info.MBID,
		})
	} else {
		queryFilters = append(queryFilters, filters.QueryFilter{
			Key: "title", Value: info.Title,
		})
	}

	existingAlbum, err := s.Get(ctx, queryFilters)
	if err == nil {
		return existingAlbum, nil
	}

	at, err := s.parseAlbumType(info.AlbumType)

	if err != nil {
		return nil, err
	}

	model := &models.Album{
		Title:         info.Title,
		Artists:       artists,
		Type:          *at,
		MusicBrainzID: info.MBID,
	}

	img, err := s.image.GetOrCreate(ctx, info.ImageURL)
	if err == nil {
		model.PictureID = &img.ID
	}

	err = s.repo.Persist(ctx, model)
	if err != nil {
		// could happen if another routine inserted the same album first
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return s.repo.FindActiveByFilters(ctx, queryFilters)
		}
		return nil, err
	}

	return model, nil
}

// Get finds a models.Album by filters
func (s *AlbumService) Get(ctx context.Context, filters []filters.QueryFilter) (*models.Album, error) {
	return s.repo.FindActiveByFilters(ctx, filters)
}

// Update partially updates a models.Album with given dto.UpdateAlbumInput
func (s *AlbumService) Update(ctx context.Context, a *models.Album, input dto.UpdateAlbumInput) (*models.Album, error) {
	updates := make(map[string]interface{})

	if input.Title != nil {
		updates["name"] = *input.Title
	}

	if input.ImageID != nil {
		updates["picture_id"] = *input.ImageID
	}

	err := s.repo.Update(ctx, a, updates)
	if err != nil {
		return nil, fmt.Errorf("unable to update album: %w", err)
	}

	return a, nil
}

func (s *AlbumService) parseAlbumType(at string) (*models.AlbumType, error) {
	m := map[models.AlbumType]struct{}{
		models.ALBUM:       {},
		models.EP:          {},
		models.SINGLE:      {},
		models.COMPILATION: {},
	}
	albumType := models.AlbumType(strings.ToUpper(at))

	_, ok := m[albumType]
	if !ok {
		return nil, fmt.Errorf("unable to parse %s as AlbumType", at)
	}

	return &albumType, nil
}
