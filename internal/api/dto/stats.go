package dto

import (
	"time"
)

type Period string

const (
	PeriodWeek    Period = "week"
	PeriodMonth   Period = "month"
	PeriodYear    Period = "year"
	PeriodOverall Period = "overall"
)

type TopArtistResult struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	PictureURL    *string `json:"picture_url"`
	ScrobbleCount *int    `json:"scrobble_count,omitempty"` // could sometimes be nil/0 when used in some queries
	// raw fields only used for building PictureURL (hidden in JSON response)
	PicturePath   string `json:"-" gorm:"column:picture_path"`
	PictureType   string `json:"-" gorm:"column:picture_type"`
	PictureDomain string `json:"-" gorm:"column:picture_domain"`
}

type TopAlbumResult struct {
	ID            int64             `json:"id"`
	Title         string            `json:"title"`
	Artists       []TopArtistResult `json:"artists" gorm:"serializer:json"`
	PictureURL    *string           `json:"picture_url"`
	ScrobbleCount *int              `json:"scrobble_count,omitempty"` // could sometimes be nil/0 when used in some queries
	// raw fields only used for building PictureURL (hidden in JSON response)
	PicturePath   string `json:"-" gorm:"column:picture_path"`
	PictureType   string `json:"-" gorm:"column:picture_type"`
	PictureDomain string `json:"-" gorm:"column:picture_domain"`
}

type TrackResult struct {
	ID            int64             `json:"id"`
	Title         string            `json:"title"`
	Artists       []TopArtistResult `json:"artists" gorm:"serializer:json"`
	PictureURL    *string           `json:"picture_url"`
	Album         TopAlbumResult    `json:"album" gorm:"serializer:json"`
	ScrobbleCount *int              `json:"scrobble_count,omitempty"` // could be nil in the history query
	// raw fields only used for building PictureURL (hidden in JSON response)
	PicturePath   string `json:"picture_path" gorm:"serializer:json"`
	PictureType   string `json:"picture_type" gorm:"serializer:json"`
	PictureDomain string `json:"picture_domain" gorm:"serializer:json"`
}

type HistoryResult struct {
	ID          int64       `json:"id"`
	ScrobbledAt time.Time   `json:"scrobbled_at"`
	Track       TrackResult `json:"track" gorm:"serializer:json"`
}
