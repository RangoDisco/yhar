package services

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
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
	baseUrl := os.Getenv("subsonicBaseUrl")
	version := os.Getenv("subsonicVersion")
	pass := os.Getenv("subsonicPassword")
	user := os.Getenv("subsonicUser")
	var nowPlaying subsonic.GetNowPlayingResponse
	url := baseUrl + "/rest/getNowPlaying?u=" + user + "&v=" + version + "&c=yhar&p=" + pass

	req, err := utils.PrepareHTTPRequest(http.MethodGet, url, "xml", &bytes.Buffer{})
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

	err = xml.NewDecoder(res.Body).Decode(&nowPlaying)
	if err != nil {
		return nil, err
	}

	return &nowPlaying, nil
}
