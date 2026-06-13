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
	filter := repositories.QueryFilter{Key: "title", Value: info.Title}
	if info.MBID != nil {
		filter = repositories.QueryFilter{
			Key: "music_brainz_id", Value: info.MBID,
		}
	}

	existingAlbum, err := s.repo.FindOneBy(ctx, []repositories.QueryFilter{filter}, "Artists.Picture", "Picture")
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
			return s.repo.FindOneBy(ctx, []repositories.QueryFilter{filter})
		}
		return nil, err
	}

	return model, nil
}

// GetByID finds a models.Album by ID
func (s *AlbumService) GetByID(ctx context.Context, id int64) (*models.Album, error) {
	return s.repo.FindOneByID(ctx, id, "Artists.Picture", "Picture")
}

// Update partially updates a models.Album with given dto.UpdateAlbumInput
func (s *AlbumService) Update(ctx context.Context, a *models.Album, input dto.UpdateAlbumInput) (*models.Album, error) {
	updates := make(map[string]interface{})

	if input.Title != nil {
		updates["title"] = *input.Title
	}

	if input.ImageID != nil {
		updates["picture_id"] = *input.ImageID
	}

	err := s.repo.Update(ctx, a.ID, updates)
	if err != nil {
		return nil, fmt.Errorf("unable to update album: %w", err)
	}

	updated, err := s.repo.FindOneByID(ctx, a.ID, "Picture")
	if err != nil {
		return nil, fmt.Errorf("unable to find updated row from db: %w", err)
	}

	return updated, nil
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
