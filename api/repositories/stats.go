package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/rangodisco/yhar/api/dto"
	"gorm.io/gorm"
)

// TODO: Refacto query as they are really similar + better handling of img

type StatsRepository struct {
	Db *gorm.DB
}

type ContentType string

const (
	Artist ContentType = "artist"
	Album  ContentType = "album"
	Track  ContentType = "track"
)

// BaseStatsQueryParams represent all params commonly used by all stats queries
type BaseStatsQueryParams struct {
	UserID string
	Start  time.Time
	End    time.Time
}

type ContentStatsQueryParams struct {
	BaseStatsQueryParams
	Interval      *string
	TargetContent ContentType
	TargetID      string
}

type PaginatedContentStatsQueryParams struct {
	ContentStatsQueryParams
	Page  int
	Limit int
}

type PaginatedStatsQueryParams struct {
	BaseStatsQueryParams
	Page  int
	Limit int
}

type StatsArtistQueryParams struct {
	PaginatedStatsQueryParams
	TrackArtistID *string
}

type StatsTrackQueryParams struct {
	PaginatedStatsQueryParams
	AlbumID       *string
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
func (r *StatsRepository) FindTopArtistsForUser(ctx context.Context, params *PaginatedContentStatsQueryParams) ([]dto.TopArtistResult, int64, error) {
	var res []dto.TopArtistResult
	var totalCount int64

	query := r.buildBaseStatQuery(ctx, params.BaseStatsQueryParams).
		Select("ar.id AS id, ar.name AS name, ar.music_brainz_id as music_brainz_id, COUNT(scrobbles.id) AS scrobble_count, "+
			"i.path AS picture_path, i.type AS picture_type, i.domain AS picture_domain").
		Where("ar.id = ?", params.TargetID).
		Joins("LEFT JOIN images i ON i.id = ar.picture_id").
		Group("ar.id, ar.name, i.path, i.type, i.domain")

	err := query.Count(&totalCount).Error
	if err != nil {
		return nil, 0, fmt.Errorf("unable to count top artists: %w", err)
	}

	err = query.Order("scrobble_count DESC").
		Scopes(Paginate(params.Page, params.Limit)).
		Find(&res).Error

	if err != nil {
		return nil, 0, fmt.Errorf("unable to find top artists: %w", err)
	}

	for i := range res {
		res[i].PictureURL = r.buildImageURL(res[i].PictureType, res[i].PictureDomain, res[i].PicturePath)
	}

	return res, totalCount, nil
}

func (r *StatsRepository) FindTopAlbumsForUser(ctx context.Context, params *PaginatedContentStatsQueryParams) ([]dto.TopAlbumResult, int64, error) {
	var res []dto.TopAlbumResult
	var totalCount int64

	query := r.buildBaseStatQuery(ctx, params.BaseStatsQueryParams).
		Select("al.id as id, al.title as title, al.music_brainz_id as music_brainz_id, " +
			"i.path AS picture_path, i.type AS picture_type, i.domain AS picture_domain, " +
			"COUNT(DISTINCT scrobbles.id) AS scrobble_count, " +
			"JSON_AGG(DISTINCT jsonb_build_object('id', ar_al.id, 'name', ar_al.name, " +
			"'picture_path', ari.path, 'picture_type', ari.type, 'picture_domain', ari.domain)) as artists").
		Joins("JOIN albums al ON al.id = tr.album_id").
		Joins("LEFT JOIN images i ON i.id = al.picture_id").
		Joins("JOIN artist_albums aral ON aral.album_id = al.id").
		Joins("JOIN artists ar_al ON ar_al.id = aral.artist_id").
		Joins("LEFT JOIN images ari ON ari.id = ar_al.picture_id").
		Group("al.id, al.title, i.path, i.type, i.domain")

	// TODO: bit weird
	if params.TargetContent == Artist {
		query = query.Where("EXISTS(SELECT 1 FROM artist_albums aral2 WHERE aral2.album_id = al.id AND aral2.artist_id = ?)", params.TargetID)
	} else {
		query = query.Where("al.id = ? ", params.TargetContent)
	}

	err := query.Count(&totalCount).Error
	if err != nil {
		return nil, 0, fmt.Errorf("unable to count top albums: %w", err)
	}

	err = query.
		Order("scrobble_count DESC").
		Scopes(Paginate(params.Page, params.Limit)).
		Find(&res).Error

	if err != nil {
		return nil, 0, fmt.Errorf("unable to find top albums: %w", err)
	}

	for i := range res {
		res[i].PictureURL = r.buildImageURL(res[i].PictureType, res[i].PictureDomain, res[i].PicturePath)
		for j := range res[i].Artists {
			a := &res[i].Artists[j]
			a.PictureURL = r.buildImageURL(a.PictureType, a.PictureDomain, a.PicturePath)
		}
	}

	return res, totalCount, nil
}

func (r *StatsRepository) FindTopTracksForUser(ctx context.Context, params *PaginatedContentStatsQueryParams) ([]dto.TrackResult, int64, error) {
	var res []dto.TrackResult
	var totalCount int64

	query := r.buildBaseStatQuery(ctx, params.BaseStatsQueryParams).
		Select("tr.id as id, tr.title as title, " +
			"i.path AS picture_path, i.type AS picture_type, i.domain AS picture_domain, " +
			"jsonb_build_object('id', al.id, 'title', al.title) as album, " +
			"COUNT(DISTINCT scrobbles.id) AS scrobble_count, " +
			"JSON_AGG(DISTINCT jsonb_build_object('id', ar.id, 'name', ar.name, " +
			"'picture_path', ari.path, 'picture_type', ari.type, 'picture_domain', ari.domain)) as artists").
		Joins("JOIN albums al ON al.id = tr.album_id").
		Joins("LEFT JOIN images i ON i.id = al.picture_id").
		Joins("LEFT JOIN images ari ON ari.id = ar.picture_id").
		Group("tr.id, tr.title, al.id, i.path, i.type, i.domain")

	if params.TargetContent == Artist {
		query = query.Where("EXISTS(SELECT 1 FROM track_artists trar2 WHERE trar2.track_id = tr.id AND trar2.artist_id = ?)", params.TargetID)
	} else if params.TargetContent == Album {
		query = query.Where("al.id = ?", params.TargetID)
	}

	err := query.Count(&totalCount).Error
	if err != nil {
		return nil, 0, fmt.Errorf("unable to count top tracks: %w", err)
	}

	err = query.Order("scrobble_count DESC").
		Scopes(Paginate(params.Page, params.Limit)).
		Find(&res).Error

	if err != nil {
		return nil, 0, fmt.Errorf("unable to find top tracks: %w", err)
	}

	for i := range res {
		res[i].PictureURL = r.buildImageURL(res[i].PictureType, res[i].PictureDomain, res[i].PicturePath)
		for j := range res[i].Artists {
			a := &res[i].Artists[j]
			a.PictureURL = r.buildImageURL(a.PictureType, a.PictureDomain, a.PicturePath)
		}
	}

	return res, totalCount, nil
}

func (r *StatsRepository) FindByUserID(ctx context.Context, params *PaginatedContentStatsQueryParams) ([]dto.HistoryResult, int64, error) {
	var res []dto.HistoryResult
	var totalCount int64

	query := r.buildBaseStatQuery(ctx, params.BaseStatsQueryParams).
		Select("scrobbles.id as id, scrobbles.scrobbled_at as scrobbled_at, json_build_object('id', tr.id, 'title', tr.title, " +
			"'picture_path', i.path, 'picture_type', i.type, 'picture_domain', i.domain,  " +
			"'album', jsonb_build_object('id', al.id, 'title', al.title), " +
			"'artists', JSON_AGG(DISTINCT jsonb_build_object('id', ar.id, 'name', ar.name, " +
			"'picture_path', ari.path, 'picture_type', ari.type, 'picture_domain', ari.domain))) as track").
		Joins("JOIN albums al ON al.id = tr.album_id").
		Joins("LEFT JOIN images i ON i.id = al.picture_id").
		Joins("LEFT JOIN images ari ON ari.id = ar.picture_id").
		Group("scrobbles.id, tr.id, tr.title, al.id, scrobbles.scrobbled_at, i.path, i.type, i.domain")

	if params.TargetContent == Artist {
		query = query.Where("EXISTS(SELECT 1 FROM track_artists trar2 WHERE trar2.track_id = tr.id AND trar2.artist_id = ?)", params.TargetID)
	}

	err := query.Count(&totalCount).Error
	if err != nil {
		return nil, 0, fmt.Errorf("unable to count scrobbles: %w", err)
	}

	err = query.Order("scrobbled_at DESC").
		Scopes(Paginate(params.Page, params.Limit)).
		Find(&res).Error

	if err != nil {
		return nil, 0, fmt.Errorf("unable to find history: %w", err)
	}

	for i := range res {
		t := &res[i].Track
		t.PictureURL = r.buildImageURL(t.PictureType, t.PictureDomain, t.PicturePath)
		for j := range res[i].Track.Artists {
			a := &res[i].Track.Artists[j]
			a.PictureURL = r.buildImageURL(a.PictureType, a.PictureDomain, a.PicturePath)
		}
	}

	return res, totalCount, nil
}

func (r *StatsRepository) FindScrobbleCountByInterval(ctx context.Context, params *ContentStatsQueryParams) ([]dto.TimelineResult, error) {
	var res []dto.TimelineResult

	subQuery := r.Db.WithContext(ctx).Select("count(*) sub_count, s.scrobbled_at::DATE as scrobble_date").
		Table("scrobbles s").
		Where("s.user_id = ? AND s.deleted_at IS null", params.UserID).
		Group("scrobble_date")
	subQuery, err := r.addModelJoin(subQuery, params.TargetContent, params.TargetID)
	if err != nil {
		return nil, fmt.Errorf("unable to build subQuery: %w", err)
	}

	query := r.Db.WithContext(ctx).Select("COALESCE(SUM(distinct s.sub_count), 0) listened_count, CASE WHEN @period = 'day' THEN ca.date::varchar WHEN @period = 'month' THEN ca.yyyymm ELSE ca.year::varchar END as listened_interval", sql.Named("period", params.Interval)).
		Table("date_calendar ca").
		Joins("LEFT JOIN (?) s ON s.scrobble_date = ca.date", subQuery).
		Group("listened_interval").
		Order("listened_interval ASC")
	query.Scopes(r.CalendarDateBetween(params.Start, params.End))

	err = query.Find(&res).Error
	if err != nil {
		return nil, fmt.Errorf("unable to fetch scrobble count: %w", err)
	}

	return res, nil
}

// addModelJoin inner join related model based on params.
func (r *StatsRepository) addModelJoin(db *gorm.DB, targetContent ContentType, targetID string) (*gorm.DB, error) {
	switch targetContent {
	case Artist:
		db.
			Joins("INNER JOIN tracks t ON s.track_id = t.id").
			Joins("INNER JOIN track_artists ta ON ta.track_id = t.id AND ta.artist_id = ?", targetID)
		break
	case Album:
		db.Joins("INNER JOIN tracks t ON s.track_id = t.id AND t.album_id = ?", targetID)
		break
	case Track:
		db.Joins("INNER JOIN tracks t ON s.track_id = t.id AND t.id = ?", targetID)
		break
	default:
		return nil, fmt.Errorf("invalid content type: %v", targetContent)
	}

	return db, nil
}

func (r *StatsRepository) buildBaseStatQuery(ctx context.Context, params BaseStatsQueryParams) *gorm.DB {
	query := r.Db.WithContext(ctx).
		Table("scrobbles").
		Joins("JOIN tracks tr ON tr.id = scrobbles.track_id").
		Joins("JOIN track_artists trar ON trar.track_id = tr.id").
		Joins("JOIN artists ar ON ar.id = trar.artist_id").
		Where("scrobbles.user_id = ?", params.UserID).
		Where("scrobbles.deleted_at IS null")

	if !params.Start.IsZero() {
		query = query.Where("scrobbles.scrobbled_at >= ?", params.Start)
	}

	if !params.End.IsZero() {
		query = query.Where("scrobbles.scrobbled_at <= ?", params.End)
	}

	return query
}

func (r *StatsRepository) CalendarDateBetween(start, end time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("date >= ?", start).
			Where("date <= ?", end)
	}
}

// buildImageURL is a shitty helper that should not exist or at be refactored
func (r *StatsRepository) buildImageURL(imageType, domain, path string) *string {
	if path == "" {
		return nil
	}

	baseURL := os.Getenv("BASE_URL")
	switch imageType {
	case "distant":
		return new(fmt.Sprintf("%s/%s", domain, path))
	default:
		return new(fmt.Sprintf("%s/images/%s", baseURL, path))
	}
}
