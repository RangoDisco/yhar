package request

type Pagination struct {
	Page  int `json:"page" form:"page" binding:"min=1"`
	Limit int `json:"limit" form:"limit" binding:"min=1,max=100"`
}
