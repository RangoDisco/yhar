package models

type AlbumImage struct {
	ID      int64 `json:"id" gorm:"primary_key"`
	AlbumID int64
	Album   Artist `gorm:"foreignKey:AlbumID;references:ID; default:null"`
	Url     string `json:"url"`
	Width   int64  `json:"width"`
	Height  int64  `json:"height"`
}
