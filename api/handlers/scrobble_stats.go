package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/api/common"
	"github.com/rangodisco/yhar/api/dto/request"
	"github.com/rangodisco/yhar/api/dto/response"
	"github.com/rangodisco/yhar/api/services"
)

type ScrobbleStatsHandler struct {
	service *services.ScrobbleStatsService
}

func NewScrobbleStatsHandler(service *services.ScrobbleStatsService) *ScrobbleStatsHandler {
	return &ScrobbleStatsHandler{service: service}
}

// GetUserTopArtists fetches the most scrobbled artists in a given period for a given user
func (h *ScrobbleStatsHandler) GetUserTopArtists(c *gin.Context) {
	h.dispatchTop[response.TopArtistResult](c, h.service.FetchUserTopArtists)
}

// GetUserTopAlbums fetches the most scrobbled albums in a given period for a given user
func (h *ScrobbleStatsHandler) GetUserTopAlbums(c *gin.Context) {
	h.dispatchTop[response.TopAlbumResult](c, h.service.FetchUserTopAlbums)
}

// GetUserTopTracks fetches the most scrobbled tracks in a given period for a given user
func (h *ScrobbleStatsHandler) GetUserTopTracks(c *gin.Context) {
	h.dispatchTop[response.TrackResult](c, h.service.FetchUserTopTracks)
}

// GetUserHistory fetches the most recent scrobbles in a given period for a given user
func (h *ScrobbleStatsHandler) GetUserHistory(c *gin.Context) {
	h.dispatchTop[response.HistoryResult](c, h.service.FetchUserHistory)
}

// GetTimeline returns a number of listens per interval over a period
func (h *ScrobbleStatsHandler) GetTimeline(c *gin.Context) {
	h.dispatchTimeline(c, nil)
}

// GetArtistTimeline returns a number of listens per interval over a period for a given artist
func (h *ScrobbleStatsHandler) GetArtistTimeline(c *gin.Context) {
	h.dispatchTimeline(c, new(request.ContentTypeArtist))
}

// GetAlbumTimeline returns a number of listens per interval over a period for a given album
func (h *ScrobbleStatsHandler) GetAlbumTimeline(c *gin.Context) {
	h.dispatchTimeline(c, new(request.ContentTypeAlbum))
}

// GetTrackTimeline returns a number of listens per interval over a period for a given track
func (h *ScrobbleStatsHandler) GetTrackTimeline(c *gin.Context) {
	h.dispatchTimeline(c, new(request.ContentTypeTrack))
}

// GetTotalListening returns the total number of scrobble and its corresponding number of ms
func (h *ScrobbleStatsHandler) GetTotalListening(c *gin.Context) {
	userID, err := common.ResolveUserID(c)
	if err != nil {
		common.RespondWithError(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	result, err := h.service.FetchTotalListeningTimeAndCount(c.Request.Context(), userID)
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, err, "Unable to fetch total listening time and count")
		return
	}

	common.RespondWithData(c, http.StatusOK, result)
}

// GetAvgDailyListens returns an avg listens per day (limit year)
func (h *ScrobbleStatsHandler) GetAvgDailyListens(c *gin.Context) {
	params, err := h.parseStatsParams(c, new(request.ContentTypeTrack))
	if err != nil {
		common.RespondWithError(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	result, err := h.service.FetchAvgDailyListens(c.Request.Context(), params)
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, err, "Unable to fetch total listening time and count")
		return
	}

	common.RespondWithData(c, http.StatusOK, result)
}

// GetDistinctListenedCount returns a number of listens for a given contentType
func (h *ScrobbleStatsHandler) GetDistinctListenedCount(c *gin.Context) {
	userID, err := common.ResolveUserID(c)
	if err != nil {
		common.RespondWithError(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	result, err := h.service.FetchDistinctListenedCount(c.Request.Context(), userID)
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, err, "Unable to fetch total listening time and count")
		return
	}

	common.RespondWithData(c, http.StatusOK, result)
}

// HELPERS //

// dispatchTop serves as a single dispatcher for all top requests, reducing boilerplate
func (h *ScrobbleStatsHandler) dispatchTop[T any](c *gin.Context, serviceMethod func(ctx context.Context, params *request.StatsParams, pagination *request.Pagination) ([]T, int64, error)) {
	params, err := h.parseStatsParams(c, nil)
	if err != nil {
		common.RespondWithError(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	pagination := parsePagination(c)

	results, total, err := serviceMethod(c.Request.Context(), params, pagination)
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, err, "Unable to fetch top items")
		return
	}

	res := common.BuildPaginatedResponse(results, pagination.Page, pagination.Limit, total)

	common.RespondWithData(c, http.StatusOK, res)
}

// dispatchTimeline serves as a single dispatcher for all timeline requests, reducing boilerplate
func (h *ScrobbleStatsHandler) dispatchTimeline(c *gin.Context, contentType *request.ContentType) {
	params, err := h.parseStatsParams(c, contentType)
	if err != nil {
		common.RespondWithError(c, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	res, err := h.service.FetchLineChartData(c.Request.Context(), params)
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, err, "Unable to fetch timeline")
		return
	}

	common.RespondWithData(c, http.StatusOK, res)
}

func parseContentParam(c *gin.Context, targetContent *request.ContentType) (string, request.ContentType, error) {
	if targetContent != nil {
		switch *targetContent {
		case request.ContentTypeArtist:
			return c.Param("artistID"), *targetContent, nil
		case request.ContentTypeAlbum:
			return c.Param("albumID"), *targetContent, nil
		case request.ContentTypeTrack:
			return c.Param("trackID"), *targetContent, nil
		default:
			return "", "", fmt.Errorf("invalid target content: %v", targetContent)
		}
	}

	if artistID := c.Query("artist"); artistID != "" {
		return artistID, request.ContentTypeArtist, nil
	}
	if trackID := c.Query("track"); trackID != "" {
		return trackID, request.ContentTypeTrack, nil
	}
	if albumID := c.Query("album"); albumID != "" {
		return albumID, request.ContentTypeAlbum, nil
	}
	return "", "", nil
}

// parseStatsParams binds all query/path params to a shared model
func (h *ScrobbleStatsHandler) parseStatsParams(c *gin.Context, initTargetContent *request.ContentType) (*request.StatsParams, error) {
	userID, err := common.ResolveUserID(c)
	if err != nil {
		return nil, err
	}

	var params request.StatsParams
	err = c.ShouldBindQuery(&params)
	if err != nil {
		return nil, err
	}

	targetID, targetContent, err := parseContentParam(c, initTargetContent)
	if err != nil {
		return nil, err
	}

	params.UserID = userID
	params.TargetContent = targetContent
	params.TargetID = targetID

	return &params, nil
}

func parsePagination(c *gin.Context) *request.Pagination {
	return &request.Pagination{
		Page:  common.ParseInt(c.Query("page"), 1),
		Limit: common.ParseInt(c.Query("limit"), 10),
	}
}
