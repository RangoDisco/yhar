package models

import "database/sql/driver"

type scrobbleOrigin string

const (
	SUBSONIC scrobbleOrigin = "SUBSONIC"
)

func (so *scrobbleOrigin) Scan(value interface{}) error {
	*so = scrobbleOrigin(value.([]byte))
	return nil
}

func (so scrobbleOrigin) Value() (driver.Value, error) {
	return string(so), nil
}

type Scrobble struct {
	ID      int64          `json:"id" gorm:"primary_key;autoIncrement"`
	Origin  scrobbleOrigin `json:"origin" gorm:"type:scrobble_origin;not null"`
	TrackID int64
	Track   Track `json:"track" gorm:"foreignkey:TrackID;references:ID;"`
	UserID  int64
	User    User `json:"user" gorm:"foreignkey:UserID;references:ID;"`
}
