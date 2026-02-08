package repositories

import (
	"time"

	"github.com/rangodisco/yhar/internal/api/config/database"
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/types/stats"
	"gorm.io/gorm"
)

type ScrobbleRepository struct {
	Db *gorm.DB
}

func NewScrobbleRepository(Db *gorm.DB) *ScrobbleRepository {
	return &ScrobbleRepository{
		Db: Db,
	}
}

func (r *ScrobbleRepository) PersistScrobble(s *models.Scrobble) error {
	res := database.GetDB().Create(&s)
	return res.Error
}

func (r *ScrobbleRepository) FindTopArtistsForUser(userID string, sd, ed time.Time, page, limit int) ([]stats.TopArtistResult, int64, error) {
	var res []stats.TopArtistResult
	var totalCount int64

	query := database.GetDB().Table("scrobbles").
		Select("ar.id AS id, ar.name AS name, i.url AS picture_url, COUNT(scrobbles.id) AS scrobble_count").
		Joins("JOIN tracks tr ON tr.id = scrobbles.track_id").
		Joins("JOIN track_artists trar ON trar.track_id = tr.id").
		Joins("JOIN artists ar ON ar.id = trar.artist_id").
		Joins("JOIN images i ON i.id = ar.picture_id").
		Where("scrobbles.user_id = ?", userID).
		Where("scrobbles.created_at >= ? AND scrobbles.created_at <= ?", sd, ed).
		Group("ar.id, i.url")

	query.Count(&totalCount)

	offset := (page - 1) * limit

	err := query.Order("scrobble_count DESC").
		Limit(limit).
		Offset(offset).
		Find(&res).Error

	return res, totalCount, err
}

func (r *ScrobbleRepository) FindTopAlbumsForUser(userID string, sd, ed time.Time, page, limit int) ([]stats.TopAlbumResult, int64, error) {
	var res []stats.TopAlbumResult
	var totalCount int64

	query := database.GetDB().Table("scrobbles").
		Select("al.id as id, al.title as title, i.url AS picture_url, COUNT(DISTINCT scrobbles.id) AS scrobble_count, JSON_AGG(DISTINCT jsonb_build_object('id', ar.id, 'name', ar.name)) as artists").
		Joins("JOIN tracks tr ON tr.id = scrobbles.track_id").
		Joins("JOIN albums al ON al.id = tr.album_id").
		Joins("JOIN images i ON i.id = al.picture_id").
		Joins("JOIN artist_albums aral ON aral.album_id = al.id").
		Joins("JOIN artists ar ON ar.id = aral.artist_id").
		Where("scrobbles.user_id = ?", userID).
		Where("scrobbles.created_at >= ? AND scrobbles.created_at <= ?", sd, ed).
		Group("al.id, al.title, i.url")

	query.Count(&totalCount)

	offset := (page - 1) * limit

	err := query.Order("scrobble_count DESC").
		Limit(limit).
		Offset(offset).
		Find(&res).Error

	return res, totalCount, err
}

func (r *ScrobbleRepository) FindTopTracksForUser(userID string, sd, ed time.Time, page, limit int) ([]stats.TopTrackResult, int64, error) {
	var res []stats.TopTrackResult
	var totalCount int64

	query := database.GetDB().Table("scrobbles").
		Select("tr.id as id, tr.title as title, jsonb_build_object('id', al.id, 'title', al.title) as album, i.url AS picture_url, COUNT(DISTINCT scrobbles.id) AS scrobble_count, JSON_AGG(DISTINCT jsonb_build_object('id', ar.id, 'name', ar.name)) as artists").
		Joins("JOIN tracks tr ON tr.id = scrobbles.track_id").
		Joins("JOIN albums al ON al.id = tr.album_id").
		Joins("JOIN images i ON i.id = al.picture_id").
		Joins("JOIN track_artists trar ON trar.track_id = tr.id").
		Joins("JOIN artists ar ON ar.id = trar.artist_id").
		Where("scrobbles.user_id = ?", userID).
		Where("scrobbles.created_at >= ? AND scrobbles.created_at <= ?", sd, ed).
		Group("tr.id, tr.title, al.id, i.url")

	query.Count(&totalCount)

	offset := (page - 1) * limit

	err := query.Order("scrobble_count DESC").
		Limit(limit).
		Offset(offset).
		Find(&res).Error

	return res, totalCount, err
}

func (r *ScrobbleRepository) FindScrobbleByUserID(userID string, page, limit int) ([]stats.ScrobbleResult, int64, error) {
	var res []stats.ScrobbleResult
	var totalCount int64

	query := database.GetDB().Table("scrobbles").
		Select("tr.id as id, tr.title as title, scrobbles.created_at as scrobbled_at, jsonb_build_object('id', al.id, 'title', al.title) as album, i.url AS picture_url, JSON_AGG(DISTINCT jsonb_build_object('id', ar.id, 'name', ar.name)) as artists").
		Joins("JOIN tracks tr ON tr.id = scrobbles.track_id").
		Joins("JOIN albums al ON al.id = tr.album_id").
		Joins("JOIN images i ON i.id = al.picture_id").
		Joins("JOIN track_artists trar ON trar.track_id = tr.id").
		Joins("JOIN artists ar ON ar.id = trar.artist_id").
		Where("scrobbles.user_id = ?", userID).
		Group("tr.id, tr.title, al.id, i.url, scrobbles.created_at")

	query.Count(&totalCount)

	offset := (page - 1) * limit

	err := query.Order("scrobbled_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&res).Error

	return res, totalCount, err
}
