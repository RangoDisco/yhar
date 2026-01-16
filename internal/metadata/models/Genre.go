package models

type Genre struct {
	ID      int64    `json:"id" gorm:"primary_key"`
	Name    string   `json:"name"`
	Artists []Artist `json:"artists" gorm:"many2many:artist_genres;"`
}
