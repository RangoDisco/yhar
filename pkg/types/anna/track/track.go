package track

type InfoByScrobbleRequest struct {
	Title  string `json:"title" binding:"required,min=2,max=255"`
	Album  string `json:"album" binding:"max=150"`
	Artist string `json:"artist" binding:"max=150"`
	Year   int64  `json:"year" binding:"gte=0,lte=9223372036854775807"`
}
