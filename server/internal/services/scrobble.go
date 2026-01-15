package services

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/rangodisco/yhar/pkg/types/anna/scrobble"
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

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(scrobbleRequest)
	if err != nil {
		return nil, err
	}

	baseUrl := os.Getenv("annaBaseUrl")
	res, err := utils.SendHTTPRequest("POST", baseUrl+"/api/tracks/by-scrobble", "json", &buf)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(res).Decode(&aRes)
	if err != nil {
		return nil, err
	}

	return &aRes, nil
}
