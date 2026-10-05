package models

type Permission struct {
	Timestamps
	ID   int64  `json:"id" gorm:"primary_key;auto_increment"`
	Name string `json:"name" gorm:"unique;not null"`
}
