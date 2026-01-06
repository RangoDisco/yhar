package models

type Artist struct {
	ID         int64   `json:"id" gorm:"primary_key"`
	ExternalID string  `json:"external_id"`
	Name       string  `json:"name"`
	Tracks     []Track `json:"tracks" gorm:"many2many:track_artists;"`
	Albums     []Album `json:"albums" gorm:"many2many:artist_albums;"`
	Genres     []Genre `json:"genres" gorm:"many2many:artist_genres;"`
}
