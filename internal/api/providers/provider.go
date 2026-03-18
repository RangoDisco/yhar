package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/time/rate"
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

var httpClient = &http.Client{Timeout: 10 * time.Second}

type MetadataProvider interface {
	Name() string
	// GetTrackByInfos fetches a track from the scrobble data
	GetTrackByInfos(ctx context.Context, infos ScrobbleData) (*TrackMetadata, error)

	// GetArtistImage fetches an artist's image URL (not all providers supports it)
	GetArtistImage(ctx context.Context, mbid, name string) (string, error)
}

// sendRequest ensures to not exceed rate limit, execute the request and unmarshal the body into the given interface
func sendRequest(ctx context.Context, url string, limiter *rate.Limiter, userAgent *string, params url.Values, result interface{}) error {
	if limiter != nil {
		err := limiter.Wait(ctx)
		if err != nil {
			return fmt.Errorf("waiting for limiter failed: %w", err)
		}
	}

	if params != nil {
		url = fmt.Sprintf("%s?%s", url, params.Encode())
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("unable to create request: %w", err)
	}

	if userAgent != nil {
		req.Header.Set("User-Agent", *userAgent)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable to perform request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("unable to close body : %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable execute request :%s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read response body: %w", err)
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		return fmt.Errorf("unable to parse JSON response: %w", err)
	}

	return nil
}
