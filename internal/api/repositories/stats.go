package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/rangodisco/yhar/internal/api/dto"
	"gorm.io/gorm"
)

// TODO: Refacto query as they are really similar

type StatsRepository struct {
	Db *gorm.DB
}

// BaseStatsQueryParams represent all params commonly used by all stats queries
type BaseStatsQueryParams struct {
	UserID string
	Start  time.Time
	End    time.Time
	Page   int
	Limit  int
}

type StatsArtistQueryParams struct {
	BaseStatsQueryParams
	TrackArtistID *string
}

type StatsAlbumQueryParams struct {
	BaseStatsQueryParams
	AlbumArtistID *string
}

type StatsTrackQueryParams struct {
	BaseStatsQueryParams
	TrackArtistID *string
	TrackID       *string
	GroupBy       *string
	OrderBy       *string
}

func NewStatsRepository(Db *gorm.DB) *StatsRepository {
	return &StatsRepository{
		Db: Db,
	}
}

// FindTopArtistsForUser finds all scrobble for given user, and group them by artist
func (r *StatsRepository) FindTopArtistsForUser(ctx context.Context, params *StatsArtistQueryParams) ([]dto.TopArtistResult, int64, error) {
	var res []dto.TopArtistResult
	var totalCount int64

	query := r.buildBaseStatQuery(ctx, params.BaseStatsQueryParams).
		Select("ar.id AS id, ar.name AS name, i.url AS picture_url, COUNT(scrobbles.id) AS scrobble_count").
		Joins("JOIN images i ON i.id = ar.picture_id").
		Group("ar.id, i.url")

	err := query.Count(&totalCount).Error
	if err != nil {
		return nil, 0, fmt.Errorf("unable to count top artists: %w", err)
	}

	err = query.Order("scrobble_count DESC").
		Limit(params.Limit).
		Offset(r.calculateOffset(params.Page, params.Limit)).
		Find(&res).Error

	if err != nil {
		return nil, 0, fmt.Errorf("unable to find top artists: %w", err)
	}

	return res, totalCount, nil
}

func (r *StatsRepository) FindTopAlbumsForUser(ctx context.Context, params *StatsAlbumQueryParams) ([]dto.TopAlbumResult, int64, error) {
	var res []dto.TopAlbumResult
	var totalCount int64

	query := r.buildBaseStatQuery(ctx, params.BaseStatsQueryParams).
		Select("al.id as id, al.title as title, i.url AS picture_url, COUNT(DISTINCT scrobbles.id) AS scrobble_count, JSON_AGG(DISTINCT jsonb_build_object('id', ar_al.id, 'name', ar_al.name, 'picture_url', ari.url)) as artists").
		Joins("JOIN albums al ON al.id = tr.album_id").
		Joins("JOIN images i ON i.id = al.picture_id").
		Joins("JOIN artist_albums aral ON aral.album_id = al.id").
		Joins("JOIN artists ar_al ON ar_al.id = aral.artist_id").
		Joins("JOIN images ari ON ari.id = ar_al.picture_id").
		Group("al.id, al.title, i.url")

	if params.AlbumArtistID != nil {
		query.Where("EXISTS(SELECT 1 FROM artist_albums aral2 WHERE aral2.album_id = al.id AND aral2.artist_id = ?)", params.AlbumArtistID)
	}

	err := query.Count(&totalCount).Error
	if err != nil {
		return nil, 0, fmt.Errorf("unable to count top albums: %w", err)
	}

	err = query.Order("scrobble_count DESC").
		Limit(params.Limit).
		Offset(r.calculateOffset(params.Page, params.Limit)).
		Find(&res).Error

	if err != nil {
		return nil, 0, fmt.Errorf("unable to find top albums: %w", err)
	}

	return res, totalCount, nil
}

// TODO: MAYBE MERGE THOSE 2
func (r *StatsRepository) FindTopTracksForUser(ctx context.Context, params *StatsTrackQueryParams) ([]dto.TrackResult, int64, error) {
	var res []dto.TrackResult
	var totalCount int64

	query := r.buildBaseStatQuery(ctx, params.BaseStatsQueryParams).
		Select("tr.id as id, tr.title as title, jsonb_build_object('id', al.id, 'title', al.title) as album, i.url AS picture_url, COUNT(DISTINCT scrobbles.id) AS scrobble_count, JSON_AGG(DISTINCT jsonb_build_object('id', ar.id, 'name', ar.name, 'picture_url', ari.url)) as artists").
		Joins("JOIN albums al ON al.id = tr.album_id").
		Joins("JOIN images i ON i.id = al.picture_id").
		Joins("JOIN images ari ON ari.id = ar.picture_id").
		Group("tr.id, tr.title, al.id, i.url")

	if params.TrackArtistID != nil {
		query.Where("EXISTS(SELECT 1 FROM track_artists trar2 WHERE trar2.track_id = tr.id AND trar2.artist_id = ?)", params.TrackArtistID)
	}

	err := query.Count(&totalCount).Error
	if err != nil {
		return nil, 0, fmt.Errorf("unable to count top tracks: %w", err)
	}

	err = query.Order("scrobble_count DESC").
		Limit(params.Limit).
		Offset(r.calculateOffset(params.Page, params.Limit)).
		Find(&res).Error

	if err != nil {
		return nil, 0, fmt.Errorf("unable to find top tracks: %w", err)
	}

	return res, totalCount, nil
}

// TODO: MAYBE MERGE THOSE 2
func (r *StatsRepository) FindScrobbleByUserID(ctx context.Context, params *StatsTrackQueryParams) ([]dto.TrackResult, int64, error) {
	var res []dto.TrackResult
	var totalCount int64

	query := r.buildBaseStatQuery(ctx, params.BaseStatsQueryParams).
		Select("tr.id as id, tr.title as title, scrobbles.scrobbled_at as scrobbled_at, jsonb_build_object('id', al.id, 'title', al.title) as album, i.url AS picture_url, JSON_AGG(DISTINCT jsonb_build_object('id', ar.id, 'name', ar.name)) as artists").
		Joins("JOIN albums al ON al.id = tr.album_id").
		Joins("JOIN images i ON i.id = al.picture_id").
		Joins("JOIN images ari ON ari.id = ar.picture_id").
		Group("tr.id, tr.title, al.id, i.url, scrobbles.scrobbled_at")

	if params.TrackArtistID != nil {
		query.Where("EXISTS(SELECT 1 FROM track_artists trar2 WHERE trar2.track_id = tr.id AND trar2.artist_id = ?)", params.TrackArtistID)
	}

	err := query.Count(&totalCount).Error
	if err != nil {
		return nil, 0, fmt.Errorf("unable to count scrobbles: %w", err)
	}

	err = query.Order("scrobbled_at DESC").
		Limit(params.Limit).
		Offset(r.calculateOffset(params.Page, params.Limit)).
		Find(&res).Error

	if err != nil {
		return nil, 0, fmt.Errorf("unable to find history: %w", err)
	}

	return res, totalCount, nil
}

func (r *StatsRepository) buildBaseStatQuery(ctx context.Context, params BaseStatsQueryParams) *gorm.DB {
	query := r.Db.WithContext(ctx).
		Table("scrobbles").
		Joins("JOIN tracks tr ON tr.id = scrobbles.track_id").
		Joins("JOIN track_artists trar ON trar.track_id = tr.id").
		Joins("JOIN artists ar ON ar.id = trar.artist_id").
		Where("scrobbles.user_id = ?", params.UserID)

	if !params.Start.IsZero() {
		query.Where("scrobbles.scrobbled_at >= ?", params.Start)
	}

	if !params.End.IsZero() {
		query.Where("scrobbles.scrobbled_at <= ?", params.End)
	}

	return query
}

func (r *StatsRepository) calculateOffset(page, limit int) int {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 1
	}
	return (page - 1) * limit
}
