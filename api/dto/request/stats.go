package request

type Period string

const (
	PeriodDay   Period = "day"
	PeriodWeek  Period = "week"
	PeriodMonth Period = "month"
	PeriodYear  Period = "year"
)

type ContentType string

// TODO: Should not stay there
const (
	ContentTypeArtist ContentType = "artist"
	ContentTypeAlbum  ContentType = "album"
	ContentTypeTrack  ContentType = "track"
)

type StatsParams struct {
	UserID        string      `json:"user_id"`
	Period        Period      `form:"period" binding:"required,oneof=day week month year overall"`
	Interval      *Period     `form:"interval" binding:"omitempty,oneof=day week month year overall"`
	TargetContent ContentType `json:"target_content,omitempty"`
	TargetID      string      `json:"target_id,omitempty"`
}
