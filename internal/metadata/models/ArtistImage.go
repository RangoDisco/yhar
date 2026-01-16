package models

type ArtistImage struct {
	ID       int64 `json:"id" gorm:"primary_key"`
	ArtistID int64
	Artist   Artist `gorm:"foreignKey:ArtistID;references:ID; default:null"`
	Url      string `json:"url"`
	Width    int64  `json:"width"`
	Height   int64  `json:"height"`
}
