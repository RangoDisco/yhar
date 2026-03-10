package models

import "time"

type Session struct {
	Timestamps
	ID          int64         `json:"id" gorm:"primary_key;autoIncrement"`
	PlayerID    string        `json:"player_id" gorm:"varchar(50)"`
	Username    string        `json:"username" gorm:"varchar(150);not null"`
	Title       string        `json:"title" gorm:"varchar(150);not null"`
	Artist      string        `json:"artist" gorm:"varchar(150)"`
	Album       string        `json:"album" gorm:"varchar(150)"`
	Duration    time.Duration `json:"duration" gorm:"not null"`
	StartedAt   time.Time     `json:"started_at" gorm:"not null"`
	LastSeenAt  time.Time     `json:"last_seen_at" gorm:"not null default:CURRENT_TIMESTAMP()"`
	CompletedAt time.Time     `json:"completed_at" gorm:"default null"`
}
