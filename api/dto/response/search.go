package response

import "github.com/rangodisco/yhar/api/dto/request"

type SearchResult struct {
	ContentType request.ContentType `json:"contentType" gorm:"column:content_type"`
	ContentID   string              `json:"contentID" gorm:"column:content_id"`
	Title       string              `json:"title" gorm:"column:title"`
	Path        string              `json:"path" gorm:"column:path"`
	Domain      string              `json:"domain" gorm:"column:domain"`
}
