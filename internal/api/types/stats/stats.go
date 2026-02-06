package stats

import "time"

type RequestPagination struct {
	Page  int
	Limit int
}

type ResponsePagination struct {
	TotalCount  int64 `json:"total_count"`
	HasNextPage bool  `json:"has_next_page"`
}

type Period string

const (
	PeriodWeek    Period = "week"
	PeriodMonth   Period = "month"
	PeriodYear    Period = "year"
	PeriodOverall Period = "overall"
)

type Params struct {
	UserID     string
	Period     Period
	Pagination RequestPagination
}

type ArtistViewModel struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	PictureURL string `json:"picture_url"`
}

type TopArtistResult struct {
	ArtistViewModel
	ScrobbleCount int `json:"scrobble_count"`
}

type AlbumViewModel struct {
	ID      int64             `json:"id"`
	Title   string            `json:"title"`
	Artists []ArtistViewModel `json:"artists" gorm:"serializer:json"`
}

type TopAlbumResult struct {
	AlbumViewModel
	PictureURL    string `json:"picture_url"`
	ScrobbleCount int    `json:"scrobble_count"`
}

type TrackViewModel struct {
	ID         int64             `json:"id"`
	Title      string            `json:"title"`
	Artists    []ArtistViewModel `json:"artists" gorm:"serializer:json"`
	PictureURL string            `json:"picture_url"`
	Album      AlbumViewModel    `json:"album" gorm:"serializer:json"`
}

type TopTrackResult struct {
	TrackViewModel
	ScrobbleCount int `json:"scrobble_count"`
}

type ScrobbleResult struct {
	TrackViewModel
	ScrobbledAt time.Time `json:"scrobbled_at"`
}

type TopResponse[T TopArtistResult | TopAlbumResult | TopTrackResult | ScrobbleResult] struct {
	Result     []T                 `json:"result"`
	Pagination *ResponsePagination `json:"pagination"`
}
