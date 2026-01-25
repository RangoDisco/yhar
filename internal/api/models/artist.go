package models

type Artist struct {
	Timestamps
	ID        int64   `json:"id" gorm:"primary_key;autoIncrement"`
	Name      string  `json:"name" gorm:"type:varchar(150);not null"`
	Genres    []Genre `json:"genres" gorm:"many2many:artist_genres;"`
	Albums    []Album `json:"albums" gorm:"many2many:artist_albums;"`
	Tracks    []Track `json:"tracks" gorm:"many2many:track_artists;"`
	PictureID int64
	Picture   Image `json:"picture" gorm:"foreignkey:PictureID"`
}
