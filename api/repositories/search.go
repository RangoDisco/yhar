package repositories

import (
	"context"
	"fmt"

	"github.com/rangodisco/yhar/api/dto/request"
	"github.com/rangodisco/yhar/api/dto/response"
	"gorm.io/gorm"
)

type SearchRepository struct {
	Db *gorm.DB
}

func NewSearchRepository(Db *gorm.DB) *SearchRepository {
	return &SearchRepository{
		Db: Db,
	}
}

func (r *SearchRepository) Search(ctx context.Context, params request.SearchQuery) ([]response.SearchResult, int64, error) {
	var res []response.SearchResult
	var totalCount int64

	query := r.Db.WithContext(ctx).
		Table("content_search").
		Select("content_type, content_id, title, path, domain, ts_rank(value, websearch_to_tsquery('english', ?)) AS rank", params.Value).
		Where("value @@ websearch_to_tsquery('english', ?)", params.Value)

	if params.ContentType != "" {
		query = query.Where("content_type = ?", params.ContentType)
	}

	err := query.Count(&totalCount).Error
	if err != nil {
		return nil, 0, fmt.Errorf("unable to search: %w", err)
	}

	err = query.Order("rank DESC").
		Scopes(Paginate(params.Pagination.Page, params.Pagination.Limit)).
		Find(&res).Error

	return res, totalCount, nil
}
