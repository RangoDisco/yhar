package dto

import "time"

type Period string

const (
	PeriodWeek    Period = "week"
	PeriodMonth   Period = "month"
	PeriodYear    Period = "year"
	PeriodOverall Period = "overall"
)

type TopArtistResult struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	PictureURL    string `json:"picture_url"`
	ScrobbleCount int    `json:"scrobble_count,omitempty"` // could sometimes be nil/0 when used in some queries
}

type TopAlbumResult struct {
	ID            int64             `json:"id"`
	Title         string            `json:"title"`
	Artists       []TopArtistResult `json:"artists" gorm:"serializer:json"`
	PictureURL    string            `json:"picture_url"`
	ScrobbleCount int               `json:"scrobble_count,omitempty"` // could sometimes be nil/0 when used in some queries
}

type TrackResult struct {
	ID            int64             `json:"id"`
	Title         string            `json:"title"`
	Artists       []TopArtistResult `json:"artists" gorm:"serializer:json"`
	PictureURL    string            `json:"picture_url"`
	Album         TopAlbumResult    `json:"album" gorm:"serializer:json"`
	ScrobbleCount int               `json:"scrobble_count,omitempty"` // could be nil in the history query
	ScrobbledAt   time.Time         `json:"scrobbled_at,omitempty"`   // could be nil in the top tracks query
}
