package providers

import (
	"context"
	"time"
)

// ScrobbleData is the data received from scrobble services
type ScrobbleData struct {
	Title  string `json:"title" binding:"required,min=2,max=255"`
	Album  string `json:"album" binding:"max=150"`
	Artist string `json:"artist" binding:"max=150"`
	Year   int64  `json:"year" binding:"gte=0,lte=9223372036854775807"`
	MBID   string `json:"mbid"`
}

type InfoResponse struct {
	Track TrackMetadata `json:"track"`
}

type TrackMetadata struct {
	Title    string           `json:"title"`
	Artists  []ArtistMetadata `json:"artists"`
	Album    AlbumMetadata    `json:"album"`
	Duration time.Duration    `json:"duration"`
	ISRC     string           `json:"isrc"`
	MBID     string           `json:"mbid"`
}

type ArtistMetadata struct {
	Name     string   `json:"title"`
	SortName string   `json:"sort_name"`
	ImageUrl string   `json:"image_url"`
	Genres   []string `json:"genres"`
	MBID     string   `json:"mbid"`
}

type AlbumMetadata struct {
	Title     string           `json:"title"`
	ImageURL  string           `json:"imageUrl"`
	Artists   []ArtistMetadata `json:"artists"`
	AlbumType string           `json:"albumType"`
	MBID      string           `json:"mbid"`
}

type MetadataProvider interface {
	Name() string
	// GetTrackByInfos fetches a track from the scrobble data
	GetTrackByInfos(ctx context.Context, infos ScrobbleData) (*TrackMetadata, error)

	// GetArtistImage fetches an artist's image URL (not all providers supports it)
	GetArtistImage(ctx context.Context, mbid, name string) (string, error)
}
