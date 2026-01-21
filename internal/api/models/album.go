package models

import "database/sql/driver"

type AlbumType string

const (
	ALBUM       AlbumType = "ALBUM"
	EP          AlbumType = "EP"
	SINGLE      AlbumType = "SINGLE"
	COMPILATION AlbumType = "COMPILATION"
)

func (at *AlbumType) Scan(value interface{}) error {
	*at = AlbumType(value.([]byte))
	return nil
}

func (at AlbumType) Value() (driver.Value, error) {
	return string(at), nil
}

type Album struct {
	Timestamps
	ID        int64     `json:"id" gorm:"primary_key;autoIncrement"`
	Title     string    `json:"title" gorm:"type:varchar(150);not null"`
	Type      AlbumType `json:"type" gorm:"type:album_type;not null"`
	Artists   []Artist  `json:"artists" gorm:"many2many:artist_albums;"`
	Genres    []Genre   `json:"genres" gorm:"many2many:album_genres;"`
	PictureID int64
	Picture   Image `json:"picture" gorm:"foreignkey:PictureID"`
}
