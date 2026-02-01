package services

import (
	"time"

	"github.com/rangodisco/yhar/internal/api/repositories"
	"github.com/rangodisco/yhar/internal/api/types/stats"
)

func BuildResponseData[T stats.TopArtistResult | stats.TopAlbumResult | stats.TopTrackResult](result []T, page, limit int, total int64) (response *stats.TopResponse[T]) {
	var hasNextPage bool

	if (int(total) / limit) >= page {
		hasNextPage = true
	} else {
		hasNextPage = false
	}

	return &stats.TopResponse[T]{
		Result: result,
		Pagination: &stats.ResponsePagination{
			TotalCount:  total,
			HasNextPage: hasNextPage,
		},
	}
}

func FetchUserTopArtists(params *stats.Params) ([]stats.TopArtistResult, int64, error) {
	sd, ed := getDateRangeFromPeriod(params.Period)
	return repositories.FindTopArtistsForUser(params.UserID, sd, ed, params.Pagination.Page, params.Pagination.Limit)
}

func FetchUserTopAlbums(params *stats.Params) ([]stats.TopAlbumResult, int64, error) {
	sd, ed := getDateRangeFromPeriod(params.Period)
	return repositories.FindTopAlbumsForUser(params.UserID, sd, ed, params.Pagination.Page, params.Pagination.Limit)
}

func FetchUserTopTracks(params *stats.Params) ([]stats.TopTrackResult, int64, error) {
	sd, ed := getDateRangeFromPeriod(params.Period)
	return repositories.FindTopTracksForUser(params.UserID, sd, ed, params.Pagination.Page, params.Pagination.Limit)
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
