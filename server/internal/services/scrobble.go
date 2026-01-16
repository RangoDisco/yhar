package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rangodisco/yhar/internal/metadata/types/scrobble"
	"github.com/rangodisco/yhar/pkg/types/subsonic"
	"github.com/rangodisco/yhar/pkg/utils"
)

// GetTrackMetadata Fetch the current playing track from all setup providers (only subsonic api for now)
// then fetches associated metadata from providers (only local db for now)
func GetTrackMetadata(entry *subsonic.Entry) (*scrobble.InfoResponse, error) {
	var aRes scrobble.InfoResponse
	scrobbleRequest := &scrobble.InfoRequest{
		Title:  entry.Title,
		Album:  entry.Album,
		Artist: entry.Artist,
	}

	aRes = anna

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(scrobbleRequest)
	if err != nil {
		return nil, err
	}

	baseUrl := os.Getenv("annaBaseUrl")
	req, err := utils.PrepareHTTPRequest("POST", baseUrl+"/api/tracks/by-scrobble", "json", &buf)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(res.Body)

	err = json.NewDecoder(res.Body).Decode(&aRes)
	if err != nil {
		return nil, err
	}

	return &aRes, nil
}
