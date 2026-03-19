package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type scrobbleOrigin string

const (
	SUBSONIC scrobbleOrigin = "SUBSONIC"
)

func (so *scrobbleOrigin) Scan(value interface{}) error {
	if value == nil {
		*so = ""
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*so = scrobbleOrigin(v)
	case string:
		*so = scrobbleOrigin(v)
	default:
		return fmt.Errorf("cannot scan type %T into scrobbleOrigin", value)
	}
	return nil
}

func (so scrobbleOrigin) Value() (driver.Value, error) {
	return string(so), nil
}

type Scrobble struct {
	Timestamps
	ID          int64          `json:"id" gorm:"primary_key;autoIncrement"`
	Origin      scrobbleOrigin `json:"origin" gorm:"type:scrobble_origin;not null"`
	TrackID     int64          `gorm:"index"`
	Track       Track          `json:"track" gorm:"foreignKey:TrackID;references:ID;"`
	UserID      int64          `gorm:"index"`
	User        User           `json:"user" gorm:"foreignKey:UserID;references:ID;"`
	ScrobbledAt time.Time      `json:"scrobbled_at" gorm:"not null"`
}
