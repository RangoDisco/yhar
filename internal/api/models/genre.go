package models

type Genre struct {
	Timestamps
	ID      int64    `json:"id" gorm:"primary_key;autoIncrement"`
	Name    string   `json:"name" gorm:"type:varchar(150);not null"`
	Artists []Artist `json:"artists" gorm:"many2many:artist_genres;"`
	Albums  []Album  `json:"albums" gorm:"many2many:album_genres;"`
}
