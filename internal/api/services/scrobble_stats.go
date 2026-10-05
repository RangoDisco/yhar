package services

import (
	"context"
	"time"

	"github.com/rangodisco/yhar/internal/api/dto"
	"github.com/rangodisco/yhar/internal/api/repositories"
)

type ScrobbleStatsService struct {
	repo *repositories.StatsRepository
}

type StatsRequest struct {
	UserID        string
	Period        dto.Period
	Start         *time.Time
	End           *time.Time
	TargetContent repositories.ContentType
	TargetID      string
}

type PaginatedStatsRequest struct {
	StatsRequest
	Pagination struct {
		Page  int
		Limit int
	}
}

func NewScrobbleStatsService(repo *repositories.StatsRepository) *ScrobbleStatsService {
	return &ScrobbleStatsService{repo: repo}
}

func (s *ScrobbleStatsService) buildBaseParams(params *PaginatedStatsRequest) *repositories.PaginatedContentStatsQueryParams {
	start, end := getDateRangeFromPeriod(params.Period)

	return &repositories.PaginatedContentStatsQueryParams{
		UserID:        params.UserID,
		Start:         start,
		End:           end,
		TargetContent: params.TargetContent,
		TargetID:      params.TargetID,
		Page:          params.Pagination.Page,
		Limit:         params.Pagination.Limit,
	}
}

func (s *ScrobbleStatsService) FetchUserTopArtists(ctx context.Context, params *PaginatedStatsRequest) ([]dto.TopArtistResult, int64, error) {
	queryParams := s.buildBaseParams(params)
	return s.repo.FindTopArtistsForUser(ctx, queryParams)
}

func (s *ScrobbleStatsService) FetchUserTopAlbums(ctx context.Context, params *PaginatedStatsRequest) ([]dto.TopAlbumResult, int64, error) {
	queryParams := s.buildBaseParams(params)
	return s.repo.FindTopAlbumsForUser(ctx, queryParams)
}

func (s *ScrobbleStatsService) FetchUserTopTracks(ctx context.Context, params *PaginatedStatsRequest) ([]dto.TrackResult, int64, error) {
	queryParams := s.buildBaseParams(params)
	return s.repo.FindTopTracksForUser(ctx, queryParams)
}

func (s *ScrobbleStatsService) FetchUserHistory(ctx context.Context, params *PaginatedStatsRequest) ([]dto.HistoryResult, int64, error) {
	queryParams := s.buildBaseParams(params)
	return s.repo.FindByUserID(ctx, queryParams)
}

func (s *ScrobbleStatsService) FetchLineChartData(ctx context.Context, params *StatsRequest) ([]dto.TimelineResult, error) {
	queryParams := &repositories.ContentStatsQueryParams{
		TargetContent: params.TargetContent,
		TargetID:      params.TargetID,
		UserID:        params.UserID,
		// TODO: fix wtf
		Interval: new(string(params.Period)),
		End:      *params.End,
		Start:    *params.Start,
	}

	return s.repo.FindScrobbleCountByInterval(ctx, queryParams)
}

func getDateRangeFromPeriod(p dto.Period) (time.Time, time.Time) {
	now := time.Now()

	switch p {
	case dto.PeriodWeek:
		return now.AddDate(0, 0, -7), now
	case dto.PeriodMonth:
		return now.AddDate(0, -1, 0), now
	case dto.PeriodYear:
		return now.AddDate(-1, 0, 0), now
	default:
		return time.Time{}, now
	}
}
