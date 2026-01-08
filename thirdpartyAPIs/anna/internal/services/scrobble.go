package services

import (
	"fmt"

	"github.com/rangodisco/yhar/pkg/types/anna/track"
	"github.com/rangodisco/yhar/thirdpartyAPIs/anna/internal/repositories"
)

func GetInfoByScrobble(scrobble track.InfoByScrobbleRequest) (*track.InfoByScrobbleResponse, error) {
	t, err := repositories.GetTrackInfoByScrobble(scrobble)
	if err != nil {
		return nil, err
	}
	fmt.Println(t)

	return &track.InfoByScrobbleResponse{}, nil
}
