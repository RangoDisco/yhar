package models

type Track struct {
	ID          int64  `json:"id" gorm:"primary_key"`
	ExternalID  string `json:"external_id"`
	Name        string `json:"name"`
	AlbumID     int64
	Album       Artist   `gorm:"foreignKey:AlbumID;references:ID; default:null"`
	TrackNumber int64    `json:"track_number"`
	DiskNumber  int64    `json:"disk_number"`
	Duration    int64    `json:"duration"`
	Artists     []Artist `json:"artists" gorm:"many2many:track_artists;"`
}
