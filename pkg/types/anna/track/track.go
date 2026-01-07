package track

import (
	"github.com/rangodisco/yhar/pkg/types/anna/album"
	"github.com/rangodisco/yhar/pkg/types/anna/artist"
)

type InfoByScrobbleRequest struct {
	Title  string `json:"title" binding:"required,min=2,max=255"`
	Album  string `json:"album" binding:"max=150"`
	Artist string `json:"artist" binding:"max=150"`
	Year   int64  `json:"year" binding:"gte=0,lte=9223372036854775807"`
}

type InfoByScrobbleResponse struct {
	ImageUrl string                          `json:"imageUrl"`
	Title    string                          `json:"title"`
	Artists  []artist.InfoByScrobbleResponse `json:"artists"`
	Album    []album.InfoByScrobbleResponse  `json:"album"`
}
