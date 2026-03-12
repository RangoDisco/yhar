package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/providers"
	"github.com/rangodisco/yhar/internal/api/repositories"
	"github.com/rangodisco/yhar/internal/api/types/filters"
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

	existingAlbum, err := s.repo.FindActiveByFilters(ctx, queryFilters)
	if err == nil {
		return existingAlbum, nil
	}

	img, _ := s.image.GetOrCreate(ctx, info.ImageURL)
	at, err := s.parseAlbumType(info.AlbumType)

	if err != nil {
		return nil, err
	}

	model := &models.Album{
		Title:         info.Title,
		Artists:       artists,
		PictureID:     &img.ID,
		Type:          *at,
		MusicBrainzID: info.MBID,
	}

	err = s.repo.Persist(ctx, model)
	if err != nil {
		return nil, err
	}
	return model, nil
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
