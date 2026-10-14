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

func (s *ScrobbleStatsService) FetchUserTopArtists(ctx context.Context, params *request.StatsParams, pagination *request.Pagination) ([]response.TopArtistResult, int64, error) {
	return s.repo.FindTopArtistsForUser(ctx, s.buildRepoParams(params), pagination.Page, pagination.Limit)
}

func (s *ScrobbleStatsService) FetchUserTopAlbums(ctx context.Context, params *request.StatsParams, pagination *request.Pagination) ([]response.TopAlbumResult, int64, error) {
	return s.repo.FindTopAlbumsForUser(ctx, s.buildRepoParams(params), pagination.Page, pagination.Limit)
}

func (s *ScrobbleStatsService) FetchUserTopTracks(ctx context.Context, params *request.StatsParams, pagination *request.Pagination) ([]response.TrackResult, int64, error) {
	return s.repo.FindTopTracksForUser(ctx, s.buildRepoParams(params), pagination.Page, pagination.Limit)
}

func (s *ScrobbleStatsService) FetchUserHistory(ctx context.Context, params *request.StatsParams, pagination *request.Pagination) ([]response.HistoryResult, int64, error) {
	return s.repo.FindByUserID(ctx, s.buildRepoParams(params), pagination.Page, pagination.Limit)
}

func (s *ScrobbleStatsService) FetchLineChartData(ctx context.Context, params *request.StatsParams) ([]response.TimelineResult, error) {
	return s.repo.FindScrobbleCountByInterval(ctx, s.buildRepoParams(params))
}

func (s *ScrobbleStatsService) FetchTotalListeningTimeAndCount(ctx context.Context, userID string) (*response.TotalListeningResult, error) {
	return s.repo.FindTotalListeningTimeAndCountByUser(ctx, userID)
}

func (s *ScrobbleStatsService) FetchAvgDailyListens(ctx context.Context, params *request.StatsParams) (*response.AvgDailyListensResult, error) {
	return s.repo.FindAvgDailyListens(ctx, s.buildRepoParams(params))
}

func (s *ScrobbleStatsService) FetchDistinctListenedCount(ctx context.Context, userID string) (*response.TotalDistinctListenedResult, error) {
	return s.repo.FindTotalListenedCount(ctx, userID)
}
func (s *ScrobbleStatsService) buildRepoParams(req *request.StatsParams) *repositories.StatsQueryParams {
	start, end := getDateRangeFromPeriod(req.Period)

	interval := req.Interval
	if interval == nil {
		interval = &req.Period
	}

	return &repositories.StatsQueryParams{
		UserID:        req.UserID,
		Start:         start,
		End:           end,
		Interval:      *interval,
		TargetContent: req.TargetContent,
		TargetID:      req.TargetID,
	}
}

func getDateRangeFromPeriod(p request.Period) (time.Time, time.Time) {
	now := time.Now()
	// Why is so tedious, am I dumb ?
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tomorrow := today.AddDate(0, 0, 1)

	switch p {
	case request.PeriodDay:
		return today.AddDate(0, 0, -1), tomorrow
	case request.PeriodWeek:
		return today.AddDate(0, 0, -7), tomorrow
	case request.PeriodMonth:
		return today.AddDate(0, -1, 0), tomorrow
	case request.PeriodYear:
		return today.AddDate(-1, 0, 0), tomorrow
	default:
		return time.Time{}, tomorrow
	}
}
