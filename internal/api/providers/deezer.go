package providers

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type DeezerProvider struct {
	name    string
	baseURL string
}

type ArtistRes struct {
	Data []Artist `json:"data"`
}

type Artist struct {
	Name    string `json:"name"`
	Picture string `json:"picture_medium"`
}

type AlbumRes struct {
	Data []Album `json:"data"`
}

type Album struct {
	Title  string `json:"title"`
	Cover  string `json:"cover_medium"`
	Artist Artist `json:"artist"`
}

func NewDeezerProvider() MetadataProvider {
	return &DeezerProvider{
		name:    "deezer",
		baseURL: "https://api.deezer.com/search",
	}
}

func (d *DeezerProvider) Name() string {
	return d.name
}

func (d *DeezerProvider) GetTrackByInfos(ctx context.Context, infos ScrobbleData) (*TrackMetadata, error) {

	return nil, fmt.Errorf("method is not implemented")
}

func (d *DeezerProvider) GetArtistImage(ctx context.Context, _, name string) (string, error) {
	params := url.Values{
		"q": {name},
	}

	endpoint := fmt.Sprintf("%s/%s", d.baseURL, "artist")

	var artistRes ArtistRes
	err := sendRequest(ctx, endpoint, nil, nil, params, &artistRes)
	if err != nil {
		return "", fmt.Errorf("unable to query artist: %w", err)
	}

	if artistRes.Data == nil || len(artistRes.Data) == 0 {
		return "", errors.New("no artist found")
	}

	artist := artistRes.Data[0]

	return artist.Picture, nil
}

func (d *DeezerProvider) GetAlbumImage(ctx context.Context, title, artist string) (string, error) {
	params := url.Values{
		"q": {title},
	}

	endpoint := fmt.Sprintf("%s/%s", d.baseURL, "album")

	var albumRes AlbumRes
	err := sendRequest(ctx, endpoint, nil, nil, params, &albumRes)
	if err != nil {
		return "", fmt.Errorf("unable to search album: %w", err)
	}

	if albumRes.Data == nil || len(albumRes.Data) == 0 {
		return "", errors.New("no album found")
	}

	for _, album := range albumRes.Data {
		if strings.ToLower(album.Title) == strings.ToLower(title) && strings.ToLower(album.Artist.Name) == strings.ToLower(artist) {
			return album.Cover, nil
		}
	}

	return "", errors.New("no matching album found")
}
