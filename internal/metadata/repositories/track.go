package repositories

import (
	"database/sql"

	"github.com/rangodisco/yhar/internal/metadata/config/database"
	"github.com/rangodisco/yhar/internal/metadata/models"
	"github.com/rangodisco/yhar/internal/metadata/types/scrobble"
)

// rawSql is god awful, but it should do the job for now, even if a bit slow sometimes.
var rawSql = `
WITH tracks_by_similarities AS (SELECT t.id                                                    as t_id,
                                       album_id                                                as tal_id,
                                       t.name                                                  as t_name,
                                       ts_rank(t.search_field, plainto_tsquery(@track_name)) as tr_rank
                                FROM tracks t
                                WHERE t.search_field @@ plainto_tsquery(@track_name)),
tracks_album AS (SELECT *,
                             al.id                                                    as al_id,
                             al.name                                                  as al_name,
                             ts_rank(al.search_field, plainto_tsquery(@album_name)) as al_rank
                      FROM tracks_by_similarities tbs
                               INNER JOIN albums al ON tbs.tal_id = al.id
                      WHERE al.search_field @@ plainto_tsquery(@album_name)),
tracks_artists AS (SELECT *,
                               ar.id                                                     as ar_id,
                               ar.name                                                   as ar_name,
                               ts_rank(ar.search_field, plainto_tsquery(@artist_name)) as ar_rank
                        FROM tracks_album tr_al
                                 INNER JOIN track_artists tr_ar ON tr_ar.track_id = tr_al.t_id
                                 INNER JOIN artists ar ON tr_ar.artist_id = ar.id
                        WHERE ar.search_field @@ plainto_tsquery(@artist_name))
SELECT *
FROM tracks_artists tar
INNER JOIN tracks t ON t.id = tar.t_id
ORDER BY tr_rank + al_rank + ar_rank DESC LIMIT 1;
`

func FindTrackInfoByScrobble(scrobble scrobble.InfoRequest) (*models.Track, error) {
	var t models.Track

	err := database.GetDB().Preload("Artists.Images").Raw(rawSql,
		sql.Named("track_name", scrobble.Title),
		sql.Named("album_name", scrobble.Album),
		sql.Named("artist_name", scrobble.Artist),
	).Find(&t).Error

	//err := database.GetDB().Preload("Artists.Images").Where("name = ?", scrobble.Title).Order("popularity desc").Last(&t).Error
	if err != nil {
		return nil, err
	}

	return &t, nil
}
