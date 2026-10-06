package request

import "time"

type Period string

const (
	PeriodWeek    Period = "week"
	PeriodMonth   Period = "month"
	PeriodYear    Period = "year"
	PeriodOverall Period = "overall"
)

type ContentType string

// TODO: Should not stay there
const (
	ContentTypeArtist ContentType = "artist"
	ContentTypeAlbum  ContentType = "album"
	ContentTypeTrack  ContentType = "track"
)

type StatsQueryParams struct {
	UserID        string      `json:"user_id"`
	Period        Period      `json:"period"`
	Start         *time.Time  `json:"start,omitempty"`
	End           *time.Time  `json:"end,omitempty"`
	TargetContent ContentType `json:"target_content,omitempty"`
	TargetID      string      `json:"target_id,omitempty"`
}

type PaginatedStatsQueryParams struct {
	StatsQueryParams
	Page  int `json:"page" binding:"min=1"`
	Limit int `json:"limit" binding:"min=1,max=100"`
}
