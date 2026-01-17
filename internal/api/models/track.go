package models

type Track struct {
	Timestamps
	ID      int64    `json:"id" gorm:"primary_key;autoIncrement"`
	Title   string   `json:"title" gorm:"type:varchar(150);not null"`
	Artists []Artist `json:"artists" gorm:"many2many:track_artists;"`
	AlbumID int64
	Album   Album `json:"album" gorm:"foreignkey:AlbumID;references:ID;default:null"`
}
