package scrobble

type InfoRequest struct {
	Title  string `json:"title" binding:"required,min=2,max=255"`
	Album  string `json:"album" binding:"max=150"`
	Artist string `json:"artist" binding:"max=150"`
	Year   int64  `json:"year" binding:"gte=0,lte=9223372036854775807"`
}

type InfoResponse struct {
	Track TrackInfo `json:"track"`
}

type TrackInfo struct {
	Title   string       `json:"title"`
	Artists []ArtistInfo `json:"artists"`
	Albums  []AlbumInfo  `json:"albums"`
}

type ArtistInfo struct {
	Name     string   `json:"title"`
	ImageUrl string   `json:"imageUrl"`
	Genres   []string `json:"genres"`
}

type AlbumInfo struct {
	Title     string       `json:"title"`
	ImageUrl  string       `json:"imageUrl"`
	Artists   []ArtistInfo `json:"artists"`
	AlbumType string       `json:"albumType"`
}
