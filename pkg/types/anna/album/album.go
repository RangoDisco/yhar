package album

import "github.com/rangodisco/yhar/pkg/types/anna/artist"

type InfoByScrobbleResponse struct {
	Title    string                          `json:"title"`
	ImageUrl string                          `json:"imageUrl"`
	Artists  []artist.InfoByScrobbleResponse `json:"artists"`
}
