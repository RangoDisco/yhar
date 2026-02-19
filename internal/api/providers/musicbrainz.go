package providers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

const (
	baseURL      = "https://musicbrainz.org/ws/2"
	coverBaseURL = "https://coverartarchive.org"

	// Used in the rate limiter to respect MusicBrainz's 1 req/s limit
	requestsPerSecond = 1
	burstSize         = 1
)

type recording struct {
	ID           string         `json:"id"`
	Title        string         `json:"title"`
	Length       int64          `json:"length"`
	ISRCs        []string       `json:"isrcs"`
	ArtistCredit []artistCredit `json:"artist-credit"`
	Releases     []release      `json:"releases"`
}

type artistCredit struct {
	Name   string `json:"name"`
	Artist struct {
		SortName string `json:"sort-name"`
		ID       string `json:"id"`
	} `json:"artist"`
}

type release struct {
	Date         string         `json:"date"`
	Title        string         `json:"title"`
	ID           string         `json:"id"`
	ArtistCredit []artistCredit `json:"artist-credit"`
	Status       string         `json:"status"`
}

type releaseWithDetails struct {
	release
	ReleaseGroup struct {
		ID          string `json:"id"`
		PrimaryType string `json:"primary-type"`
	} `json:"release-group"`
	CoverURL string `json:"cover-url"`
}

type coverArtResponse struct {
	Images  []coverArtImage `json:"images"`
	Release string          `json:"release"`
}

type coverArtImage struct {
	Front bool   `json:"front"`
	Image string `json:"image"`
}

type MusicBrainzProvider struct {
	name         string
	baseURL      string
	coverBaseURL string
	userAgent    string
	client       *http.Client
	limiter      *rate.Limiter
}

func NewMusicBrainzProvider() MetadataProvider {
	version := os.Getenv("APP_VERSION")
	return &MusicBrainzProvider{
		name:         "musicbrainz",
		baseURL:      baseURL,
		coverBaseURL: coverBaseURL,
		// TODO: change
		userAgent: "yhar/" + version,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		limiter: rate.NewLimiter(rate.Every(requestsPerSecond*time.Second), burstSize),
	}
}

func (p *MusicBrainzProvider) Name() string {
	return p.name
}

func (p *MusicBrainzProvider) GetTrackByInfos(ctx context.Context, data ScrobbleData) (*TrackMetadata, error) {
	if data.MBID == "" && data.Title == "" {
		return nil, errors.New("no MBID nor title was provided")
	}

	var rec *recording
	var err error

	// If MBID is provided, use it
	if data.MBID != "" {
		rec, err = p.getTrackByMBID(ctx, data.MBID)
	} else {
		// Otherwise, search by title and artist
		rec, err = p.searchTrack(ctx, data.Title, data.Artist)
	}

	if err != nil {
		return nil, fmt.Errorf("unable to search track: %w", err)
	}

	track := p.convertRecordingToTrack(ctx, rec)

	// Find the "best" release and put it as track
	album, err := p.findBestRelease(ctx, rec.Releases)
	if err != nil {
		return nil, fmt.Errorf("unable to find best release: %w", err)
	}

	track.Album = *album

	return track, nil
}

func (p *MusicBrainzProvider) getTrackByMBID(ctx context.Context, mbid string) (*recording, error) {
	endpoint := fmt.Sprintf("%s/recording/%s", p.baseURL, mbid)
	params := url.Values{
		"fmt": {"json"},
	}
	params.Set("inc", strings.Join([]string{"artist-credits", "releases"}, "+"))

	var recordingRes recording
	err := p.sendRequest(ctx, endpoint, params, &recordingRes)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch recording by ID : %w", err)
	}

	return &recordingRes, nil
}

func (p *MusicBrainzProvider) searchTrack(ctx context.Context, title, artist string) (*recording, error) {
	endpoint := fmt.Sprintf("%s/recording", p.baseURL)
	query := fmt.Sprintf("recording:%s", title)

	if artist != "" {
		query = query + fmt.Sprintf(" AND artistname=%s", artist)
	}

	// Artist-credits and releases are already inc for some reason
	params := url.Values{
		"fmt":   {"json"},
		"query": {url.QueryEscape(title)},
	}

	var recordingRes []recording
	err := p.sendRequest(ctx, endpoint, params, &recordingRes)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch recording by name : %w", err)
	}

	// TODO: determine best recording in case multiple were returned
	return &recordingRes[0], nil
}

// GetArtistImage is not implemented by this provider as MusicBrainz doesn't provide images for artists
func (p *MusicBrainzProvider) GetArtistImage(ctx context.Context, mbid, name string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

// convertRecordingToTrack converts a MusicBrainz's recording response to a standardized TrackMetadata Provider response
func (p *MusicBrainzProvider) convertRecordingToTrack(ctx context.Context, rec *recording) *TrackMetadata {
	track := &TrackMetadata{
		Title:    rec.Title,
		MBID:     rec.ID,
		Duration: time.Duration(rec.Length) * time.Millisecond,
	}

	// Convert artists
	var artists []ArtistMetadata
	for _, credit := range rec.ArtistCredit {
		artists = append(artists, ArtistMetadata{
			Name:     credit.Name,
			SortName: credit.Artist.SortName,
			MBID:     credit.Artist.ID,
		})
	}
	track.Artists = artists

	return track
}

// convertRecordingToAlbum converts MusicBrainz's release response to a standardized AlbumMetadata Provider response
func (p *MusicBrainzProvider) convertReleaseToAlbum(ctx context.Context, releaseDetails *releaseWithDetails) *AlbumMetadata {
	album := &AlbumMetadata{
		Title:     releaseDetails.Title,
		MBID:      releaseDetails.ID,
		AlbumType: releaseDetails.ReleaseGroup.PrimaryType,
		ImageURL:  releaseDetails.CoverURL,
	}

	// Convert Artists
	var artists []ArtistMetadata
	for _, credit := range releaseDetails.ArtistCredit {
		artists = append(artists, ArtistMetadata{
			Name:     credit.Name,
			SortName: credit.Artist.SortName,
			MBID:     credit.Artist.ID,
		})
	}

	album.Artists = artists

	return album
}

// findBestRelease selects the "best" release and fetch its complete metadata
func (p *MusicBrainzProvider) findBestRelease(ctx context.Context, releases []release) (*AlbumMetadata, error) {
	if len(releases) == 0 {
		return nil, fmt.Errorf("no releases were found")
	}

	// I have no fucking clue on how to determine the best release, so I just take the first official + dated one
	var bestRelease *release
	for _, r := range releases {
		if r.Date != "" && r.Status == "Official" {
			bestRelease = &r
			break
		}
	}

	// Default to the first release
	if bestRelease == nil {
		bestRelease = &releases[0]
	}

	// Including release doesn't include its cover and type for some reason ? So we have to fetch it
	details, err := p.getReleaseDetails(ctx, bestRelease.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to get album's cover: %w", err)
	}

	album := p.convertReleaseToAlbum(ctx, details)

	return album, nil
}

// getReleaseDetails fetches complete release information from MusicBrainz.
func (p *MusicBrainzProvider) getReleaseDetails(ctx context.Context, releaseID string) (*releaseWithDetails, error) {
	endpoint := fmt.Sprintf("%s/release/%s", p.baseURL, releaseID)
	params := url.Values{
		"fmt": {"json"},
	}
	params.Set("inc", strings.Join([]string{"artists", "artist-credits", "release-groups"}, "+"))

	var details releaseWithDetails
	err := p.sendRequest(ctx, endpoint, params, &details)
	if err != nil {
		return nil, fmt.Errorf("fetch release details: %w", err)
	}

	// Also fetch the album cover
	cover, err := p.getAlbumCover(ctx, releaseID)
	if err != nil {
		return nil, err
	}

	details.CoverURL = cover

	return &details, nil
}

// getAlbumCover fetches the album cover URL from Cover Art Archive.
// 307 redirect to an index.json file, if there is a release with this MBID.
// 400 if {mbid} cannot be parsed as a valid UUID.
// 404 if there is no release with this MBID.
// 405 if the request method is not one of GET or HEAD.
// 406 if the server is unable to generate a response suitable to the Accept header.
// 503 if the user has exceeded their rate limit.
func (p *MusicBrainzProvider) getAlbumCover(ctx context.Context, releaseID string) (string, error) {
	endpoint := fmt.Sprintf("%s/release/%s", p.coverBaseURL, releaseID)

	var coverResp coverArtResponse
	err := p.sendRequest(ctx, endpoint, nil, &coverResp)
	if err != nil {
		return "", nil // No cover art available, not an error
	}

	// TODO: Handle multiple images
	// Find front cover
	for _, image := range coverResp.Images {
		if image.Front {
			return image.Image, nil
		}
	}

	// Fallback to first image if no front cover marked
	if len(coverResp.Images) > 0 {
		return coverResp.Images[0].Image, nil
	}

	return "", nil
}

// sendRequest ensure to not exceed MusicBrainz's rate limit (1 req/s), execute the request and unmarshal the body into the given interface
func (p *MusicBrainzProvider) sendRequest(ctx context.Context, url string, params url.Values, result interface{}) error {
	err := p.limiter.Wait(ctx)
	if err != nil {
		return fmt.Errorf("waiting for limiter failed: %w", err)
	}

	if params != nil {
		url = fmt.Sprintf("%s?%s", url, params.Encode())
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Set("User-Agent", p.userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := p.client.Do(req)
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
