package stats

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

type TopArtistResult struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	PictureURL    string `json:"picture_url"`
	ScrobbleCount int    `json:"scrobble_count"`
}

type TopAlbumResult struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Artists []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"artists" gorm:"serializer:json"`

	PictureURL    string `json:"picture_url"`
	ScrobbleCount int    `json:"scrobble_count"`
}

type TopTrackResult struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Artists []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"artists" gorm:"serializer:json"`
	PictureURL    string `json:"picture_url"`
	ScrobbleCount int    `json:"scrobble_count"`
	Album         struct {
		ID    int64  `json:"id"`
		Title string `json:"title"`
	} `json:"album" gorm:"serializer:json"`
}

type TopResponse[T TopArtistResult | TopAlbumResult | TopTrackResult] struct {
	Result     []T                 `json:"result"`
	Pagination *ResponsePagination `json:"pagination"`
}
