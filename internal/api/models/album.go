package models

import "database/sql/driver"

type albumType string

const (
	ALBUM       albumType = "ALBUM"
	EP          albumType = "EP"
	SINGLE      albumType = "SINGLE"
	COMPILATION albumType = "COMPILATION"
)

func (at *albumType) Scan(value interface{}) error {
	*at = albumType(value.([]byte))
	return nil
}

func (at albumType) Value() (driver.Value, error) {
	return string(at), nil
}

type Album struct {
	Timestamps
	ID      int64     `json:"id" gorm:"primary_key;autoIncrement"`
	Title   string    `json:"title" gorm:"type:varchar(150);not null"`
	Type    albumType `json:"type" gorm:"type:album_type;not null"`
	Artists []Artist  `json:"artists" gorm:"many2many:artist_albums;"`
	Genres  []Genre   `json:"genres" gorm:"many2many:album_genres;"`
}
