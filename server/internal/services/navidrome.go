package services

import (
	"encoding/xml"
	"net/http"
	"os"

	"github.com/rangodisco/yhar/pkg/types/subsonic"
	"github.com/rangodisco/yhar/pkg/utils"
)

var (
	baseUrl = os.Getenv("subsonicBaseUrl")
	version = os.Getenv("subsonicVersion")
	pass    = os.Getenv("subsonicPassword")
	user    = os.Getenv("subsonicUser")
)

// GetNowPlaying fetch current playing tracks from all sources (only subsonic for now)
func GetNowPlaying() (*subsonic.GetNowPlayingResponse, error) {
	var nowPlaying subsonic.GetNowPlayingResponse
	url := baseUrl + "/rest/getNowPlaying?u=" + user + "&v=" + version + "&c=yhar&p=" + pass

	res, err := utils.SendHTTPRequest(http.MethodGet, url, "xml", nil)
	if err != nil {
		return nil, err
	}

	err = xml.NewDecoder(res).Decode(&nowPlaying)
	if err != nil {
		return nil, err
	}

	return &nowPlaying, nil
}
