package response

import (
	"time"
)

type TopArtistResult struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	MusicBrainzID *string `json:"music_brainz_id,omitempty"`
	ScrobbleCount *int    `json:"scrobble_count,omitempty"` // could sometimes be nil/0 when used in some queries
	PicturePath   string  `json:"picturePath" gorm:"column:picturePath"`
	PictureDomain string  `json:"pictureDomain" gorm:"column:pictureDomain"`
}

type TopAlbumResult struct {
	ID            int64             `json:"id"`
	Title         string            `json:"title"`
	Artists       []TopArtistResult `json:"artists" gorm:"serializer:json"`
	MusicBrainzID *string           `json:"music_brainz_id,omitempty"`
	ScrobbleCount *int              `json:"scrobble_count,omitempty"` // could sometimes be nil/0 when used in some queries
	PicturePath   string            `json:"picturePath" gorm:"column:picturePath"`
	PictureDomain string            `json:"pictureDomain" gorm:"column:pictureDomain"`
}

type TrackResult struct {
	ID            int64             `json:"id"`
	Title         string            `json:"title"`
	Artists       []TopArtistResult `json:"artists" gorm:"serializer:json"`
	Album         TopAlbumResult    `json:"album" gorm:"serializer:json"`
	ScrobbleCount *int              `json:"scrobble_count,omitempty"` // could be nil in the history query
	PicturePath   string            `json:"picturePath" gorm:"column:picturePath"`
	PictureDomain string            `json:"pictureDomain" gorm:"column:pictureDomain"`
}

type HistoryResult struct {
	ID          int64       `json:"id"`
	ScrobbledAt time.Time   `json:"scrobbled_at"`
	Track       TrackResult `json:"track" gorm:"serializer:json"`
}

type TimelineResult struct {
	ListeningCount int64  `json:"listenedCount" gorm:"column:listened_count"`
	PeriodInterval string `json:"periodInterval" gorm:"column:listened_interval"`
}

type TotalListeningResult struct {
	TotalScrobbles     int64 `json:"totalScrobbles" gorm:"column:total_scrobbles"`
	TotalListeningTime int64 `json:"totalDuration" gorm:"column:total_duration"`
}

type AvgDailyListensResult struct {
	AvgCount float64 `json:"avgCount" gorm:"column:avg_count"`
}

type TotalDistinctListenedResult struct {
	TotalCount int64 `json:"totalCount" gorm:"column:total_count"`
}
