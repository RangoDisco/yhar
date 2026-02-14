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

type StatsRequestParams struct {
	UserID     string
	Period     dto.Period
	ArtistID   *string
	TrackID    *string
	Pagination struct {
		Page  int
		Limit int
	}
}

func NewScrobbleStatsService(repo *repositories.StatsRepository) *ScrobbleStatsService {
	return &ScrobbleStatsService{repo: repo}
}

func (s *ScrobbleStatsService) buildBaseParams(params *StatsRequestParams) repositories.BaseStatsQueryParams {
	start, end := getDateRangeFromPeriod(params.Period)

	return repositories.BaseStatsQueryParams{
		UserID: params.UserID,
		Start:  start,
		End:    end,
		Page:   params.Pagination.Page,
		Limit:  params.Pagination.Limit,
	}
}

func (s *ScrobbleStatsService) FetchUserTopArtists(ctx context.Context, params *StatsRequestParams) ([]dto.TopArtistResult, int64, error) {
	queryParams := &repositories.StatsArtistQueryParams{
		BaseStatsQueryParams: s.buildBaseParams(params),
		TrackArtistID:        params.ArtistID,
	}
	return s.repo.FindTopArtistsForUser(ctx, queryParams)
}

func (s *ScrobbleStatsService) FetchUserTopAlbums(ctx context.Context, params *StatsRequestParams) ([]dto.TopAlbumResult, int64, error) {
	queryParams := &repositories.StatsAlbumQueryParams{
		BaseStatsQueryParams: s.buildBaseParams(params),
		AlbumArtistID:        params.ArtistID,
	}
	return s.repo.FindTopAlbumsForUser(ctx, queryParams)
}

func (s *ScrobbleStatsService) FetchUserTopTracks(ctx context.Context, params *StatsRequestParams) ([]dto.TrackResult, int64, error) {
	queryParams := &repositories.StatsTrackQueryParams{
		BaseStatsQueryParams: s.buildBaseParams(params),
		TrackID:              params.TrackID,
		TrackArtistID:        params.ArtistID,
	}
	return s.repo.FindTopTracksForUser(ctx, queryParams)
}

func (s *ScrobbleStatsService) FetchUserHistory(ctx context.Context, params *StatsRequestParams) ([]dto.TrackResult, int64, error) {
	queryParams := &repositories.StatsTrackQueryParams{
		BaseStatsQueryParams: s.buildBaseParams(params),
		TrackID:              params.TrackID,
		TrackArtistID:        params.ArtistID,
	}
	return s.repo.FindScrobbleByUserID(ctx, queryParams)
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
