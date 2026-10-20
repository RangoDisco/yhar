package request

type SearchQuery struct {
	Value       string      `json:"value" form:"value" binding:"required"`
	ContentType ContentType `json:"contentType" form:"type" binding:"omitempty,oneof=artists albums tracks"`
	Pagination  Pagination  `json:"pagination"`
}
