package services

import (
	"time"

	"github.com/rangodisco/yhar/internal/api/repositories"
	"github.com/rangodisco/yhar/internal/api/types/stats"
)

type ScrobbleStatsService struct {
	sRepo *repositories.ScrobbleRepository
}

func NewScrobbleStatsService(sRepo *repositories.ScrobbleRepository) *ScrobbleStatsService {
	return &ScrobbleStatsService{sRepo: sRepo}
}

func (s *ScrobbleStatsService) BuildResponseData(result interface{}, page, limit int, total int64) interface{} {
	var hasNextPage bool
	if (int(total) / limit) >= page {
		hasNextPage = true
	} else {
		hasNextPage = false
	}

	pagination := &stats.ResponsePagination{
		TotalCount:  total,
		HasNextPage: hasNextPage,
	}

	switch v := result.(type) {
	case []stats.TopArtistResult:
		return &stats.TopResponse[stats.TopArtistResult]{
			Result:     v,
			Pagination: pagination,
		}
	case []stats.TopAlbumResult:
		return &stats.TopResponse[stats.TopAlbumResult]{
			Result:     v,
			Pagination: pagination,
		}
	case []stats.TopTrackResult:
		return &stats.TopResponse[stats.TopTrackResult]{
			Result:     v,
			Pagination: pagination,
		}
	case []stats.ScrobbleResult:
		return &stats.TopResponse[stats.ScrobbleResult]{
			Result:     v,
			Pagination: pagination,
		}
	default:
		return nil
	}
}

func (s *ScrobbleStatsService) FetchUserTopArtists(params *stats.Params) ([]stats.TopArtistResult, int64, error) {
	sd, ed := getDateRangeFromPeriod(params.Period)
	return s.sRepo.FindTopArtistsForUser(params.UserID, sd, ed, params.Pagination.Page, params.Pagination.Limit)
}

func (s *ScrobbleStatsService) FetchUserTopAlbums(params *stats.Params) ([]stats.TopAlbumResult, int64, error) {
	sd, ed := getDateRangeFromPeriod(params.Period)
	return s.sRepo.FindTopAlbumsForUser(params.UserID, sd, ed, params.Pagination.Page, params.Pagination.Limit)
}

func (s *ScrobbleStatsService) FetchUserTopTracks(params *stats.Params) ([]stats.TopTrackResult, int64, error) {
	sd, ed := getDateRangeFromPeriod(params.Period)
	return s.sRepo.FindTopTracksForUser(params.UserID, sd, ed, params.Pagination.Page, params.Pagination.Limit)
}

func (s *ScrobbleStatsService) FetchUserHistory(params *stats.Params) ([]stats.ScrobbleResult, int64, error) {
	return s.sRepo.FindScrobbleByUserID(params.UserID, params.Pagination.Page, params.Pagination.Limit)
}

func getDateRangeFromPeriod(p stats.Period) (time.Time, time.Time) {
	now := time.Now()
	switch p {
	case stats.PeriodWeek:
		return now.AddDate(0, 0, -7), now
	case stats.PeriodMonth:
		return now.AddDate(0, -1, 0), now
	case stats.PeriodYear:
		return now.AddDate(-1, 0, 0), now
	default:
		return time.Time{}, now
	}
}
