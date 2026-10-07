package services

import (
	"context"
	"time"

	"github.com/rangodisco/yhar/api/dto/request"
	"github.com/rangodisco/yhar/api/dto/response"
	"github.com/rangodisco/yhar/api/repositories"
)

type ScrobbleStatsService struct {
	repo *repositories.StatsRepository
}

func NewScrobbleStatsService(repo *repositories.StatsRepository) *ScrobbleStatsService {
	return &ScrobbleStatsService{repo: repo}
}

func (s *ScrobbleStatsService) buildRepoParams(req *request.StatsQueryParams) *repositories.StatsQueryParams {
	start, end := getDateRangeFromPeriod(req.Period)

	return &repositories.StatsQueryParams{
		UserID:        req.UserID,
		Start:         start,
		End:           end,
		Interval:      req.Interval,
		TargetContent: req.TargetContent,
		TargetID:      req.TargetID,
	}
}

func (s *ScrobbleStatsService) buildPaginatedRepoParams(req *request.PaginatedStatsQueryParams) *repositories.PaginatedStatsQueryParams {
	return &repositories.PaginatedStatsQueryParams{
		StatsQueryParams: *s.buildRepoParams(&req.StatsQueryParams),
		Page:             req.Page,
		Limit:            req.Limit,
	}
}

func (s *ScrobbleStatsService) FetchUserTopArtists(ctx context.Context, params *request.PaginatedStatsQueryParams) ([]response.TopArtistResult, int64, error) {
	return s.repo.FindTopArtistsForUser(ctx, s.buildPaginatedRepoParams(params))
}

func (s *ScrobbleStatsService) FetchUserTopAlbums(ctx context.Context, params *request.PaginatedStatsQueryParams) ([]response.TopAlbumResult, int64, error) {
	return s.repo.FindTopAlbumsForUser(ctx, s.buildPaginatedRepoParams(params))
}

func (s *ScrobbleStatsService) FetchUserTopTracks(ctx context.Context, params *request.PaginatedStatsQueryParams) ([]response.TrackResult, int64, error) {
	return s.repo.FindTopTracksForUser(ctx, s.buildPaginatedRepoParams(params))
}

func (s *ScrobbleStatsService) FetchUserHistory(ctx context.Context, params *request.PaginatedStatsQueryParams) ([]response.HistoryResult, int64, error) {
	return s.repo.FindByUserID(ctx, s.buildPaginatedRepoParams(params))
}

func (s *ScrobbleStatsService) FetchLineChartData(ctx context.Context, params *request.StatsQueryParams) ([]response.TimelineResult, error) {
	return s.repo.FindScrobbleCountByInterval(ctx, s.buildRepoParams(params))
}

func getDateRangeFromPeriod(p request.Period) (time.Time, time.Time) {
	now := time.Now()

	switch p {
	case request.PeriodDay:
		return now.AddDate(0, 0, -1), now
	case request.PeriodWeek:
		return now.AddDate(0, 0, -7), now
	case request.PeriodMonth:
		return now.AddDate(0, -1, 0), now
	case request.PeriodYear:
		return now.AddDate(-1, 0, 0), now
	default:
		return time.Time{}, now
	}
}
