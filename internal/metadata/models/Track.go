package models

type Track struct {
	ID          int64  `json:"id" gorm:"primary_key"`
	ExternalID  string `json:"external_id"`
	Name        string `json:"name"`
	AlbumID     int64
	Album       Album    `gorm:"foreignKey:AlbumID;references:ID; default:null"`
	TrackNumber int64    `json:"track_number"`
	DiscNumber  int64    `json:"disc_number"`
	Duration    int64    `json:"duration"`
	Popularity  int64    `json:"popularity"`
	Artists     []Artist `json:"artists" gorm:"many2many:track_artists;"`
}
