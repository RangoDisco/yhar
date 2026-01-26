package models

type User struct {
	Timestamps
	ID         int64      `json:"id" gorm:"primary_key;autoIncrement"`
	Username   string     `json:"username" gorm:"type:varchar(75);not null"`
	Origin     string     `json:"origin" gorm:"type:varchar(255);not null"`
	ExternalID string     `json:"external_id" gorm:"type:varchar(255);default: null;uniqueIndex"`
	Scrobbles  []Scrobble `json:"scrobbles" gorm:"foreignkey:UserID;"`
	Password   string     `json:"password" gorm:"type:varchar(255);not null"`
	RoleID     int64      `json:"role_id" gorm:"not null"`
	Role
}
