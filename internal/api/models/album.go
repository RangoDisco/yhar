package models

import (
	"database/sql/driver"
	"fmt"
)

type AlbumType string

const (
	ALBUM       AlbumType = "ALBUM"
	EP          AlbumType = "EP"
	SINGLE      AlbumType = "SINGLE"
	COMPILATION AlbumType = "COMPILATION"
)

func (at *AlbumType) Scan(value interface{}) error {
	if value == nil {
		*at = ""
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*at = AlbumType(v)
		return nil
	case string:
		*at = AlbumType(v)
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into AlbumType", value)
	}
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
