package models

type Image struct {
	Timestamps
	ID  int64  `json:"id" gorm:"primary_key;autoIncrement"`
	Url string `json:"url" gorm:"not null"`
}
