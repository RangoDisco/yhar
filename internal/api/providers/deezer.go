package providers

import (
	"context"
	"fmt"
	"net/url"
)

type DeezerProvider struct {
	name    string
	baseURL string
}

type ArtistRes struct {
	Data []Artist `json:"data"`
}

type Artist struct {
	Picture string `json:"picture_medium"`
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
		return "", fmt.Errorf("no artist found")
	}

	artist := artistRes.Data[0]

	return artist.Picture, nil
}
