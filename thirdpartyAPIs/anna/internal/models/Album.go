package models

type Album struct {
	ID         int64  `json:"id" gorm:"primary_key"`
	ExternalID string `json:"external_id"`
	Name       string `json:"name"`
	// TODO: Enum
	AlbumType   string       `json:"album_type"`
	ReleaseDate string       `json:"release_date"`
	TotalTracks int64        `json:"total_tracks"`
	Artists     []Artist     `json:"artists" gorm:"many2many:artist_albums;"`
	Images      []AlbumImage `json:"images" gorm:"foreignKey:AlbumID;references:ID"`
}
