package request

type Period string

const (
	PeriodDay   Period = "day"
	PeriodWeek  Period = "week"
	PeriodMonth Period = "month"
	PeriodYear  Period = "year"
)

type ContentType string

const (
	ContentTypeArtist ContentType = "artists"
	ContentTypeAlbum  ContentType = "albums"
	ContentTypeTrack  ContentType = "tracks"
)

type StatsParams struct {
	UserID        string      `json:"user_id"`
	Period        Period      `form:"period" binding:"required,oneof=day week month year overall"`
	Interval      *Period     `form:"interval" binding:"omitempty,oneof=day week month year overall"`
	TargetContent ContentType `json:"target_content,omitempty"`
	TargetID      string      `json:"target_id,omitempty"`
}
