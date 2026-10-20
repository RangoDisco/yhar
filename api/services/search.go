package services

import (
	"context"

	"github.com/rangodisco/yhar/api/dto/request"
	"github.com/rangodisco/yhar/api/dto/response"
	"github.com/rangodisco/yhar/api/repositories"
)

type SearchService struct {
	repo *repositories.SearchRepository
}

func NewSearchService(repo *repositories.SearchRepository) *SearchService {
	return &SearchService{repo: repo}
}

func (s *SearchService) Search(ctx context.Context, query request.SearchQuery) ([]response.SearchResult, int64, error) {
	return s.repo.Search(ctx, query)
}
