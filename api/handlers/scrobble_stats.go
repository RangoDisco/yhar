package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/api/common"
	"github.com/rangodisco/yhar/api/dto/request"
	"github.com/rangodisco/yhar/api/services"
)

// TODO: should be refactored, too much repetition + god awful param handling
type ScrobbleStatsHandler struct {
	service *services.ScrobbleStatsService
}

func NewScrobbleStatsHandler(service *services.ScrobbleStatsService) *ScrobbleStatsHandler {
	return &ScrobbleStatsHandler{service: service}
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

func (h *ScrobbleStatsHandler) parseStatsParams(c *gin.Context, initTargetContent *request.ContentType) (*request.StatsQueryParams, error) {
	userID, err := common.ResolveUserID(c)
	if err != nil {
		return nil, err
	}

	targetID, targetContent, err := parseContentParam(c, initTargetContent)
	if err != nil {
		return nil, err
	}

	return &request.StatsQueryParams{
		UserID:        userID,
		Period:        request.Period(c.DefaultQuery("period", string(request.PeriodWeek))),
		Interval:      request.Period(c.DefaultQuery("interval", string(request.PeriodDay))),
		TargetContent: targetContent,
		TargetID:      targetID,
	}, nil
}

func (h *ScrobbleStatsHandler) parsePaginatedStatsParams(c *gin.Context) (*request.PaginatedStatsQueryParams, error) {
	userID, err := common.ResolveUserID(c)
	if err != nil {
		return nil, err
	}

	targetID, targetContent, err := parseContentParam(c, nil)
	if err != nil {
		return nil, err
	}

	// Parse and validate pagination
	page := common.ParseInt(c.Query("page"), 1)
	limit := common.ParseInt(c.Query("limit"), 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 1
	}
	if limit > 100 {
		limit = 100
	}

	// Build params
	return &request.PaginatedStatsQueryParams{
		UserID:        userID,
		Period:        request.Period(c.DefaultQuery("period", string(request.PeriodWeek))),
		TargetContent: targetContent,
		TargetID:      targetID,
		Page:          page,
		Limit:         limit,
	}, nil
}

// GetUserTopArtists fetches the most scrobbled artists in a given period for a given user
func (h *ScrobbleStatsHandler) GetUserTopArtists(c *gin.Context) {
	params, err := h.parsePaginatedStatsParams(c)
	if err != nil {
		common.RespondWithError(c, http.StatusBadRequest, err, "Invalid body")
		return
	}

	results, total, err := h.service.FetchUserTopArtists(c.Request.Context(), params)
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, err, "Unable to fetch top artists")
		return
	}

	res := common.BuildPaginatedResponse(results, params.Page, params.Limit, total)

	common.RespondWithData(c, http.StatusOK, res)
}

// GetUserTopAlbums fetches the most scrobbled albums in a given period for a given user
func (h *ScrobbleStatsHandler) GetUserTopAlbums(c *gin.Context) {
	params, err := h.parsePaginatedStatsParams(c)
	if err != nil {
		common.RespondWithError(c, http.StatusBadRequest, err, "Invalid body")
		return
	}

	results, total, err := h.service.FetchUserTopAlbums(c.Request.Context(), params)
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, err, "Unable to fetch top albums")
		return
	}

	res := common.BuildPaginatedResponse(results, params.Page, params.Limit, total)

	common.RespondWithData(c, http.StatusOK, res)
}

// GetUserTopTracks fetches the most scrobbled tracks in a given period for a given user
func (h *ScrobbleStatsHandler) GetUserTopTracks(c *gin.Context) {
	params, err := h.parsePaginatedStatsParams(c)
	if err != nil {
		common.RespondWithError(c, http.StatusBadRequest, err, "Invalid body")
		return
	}

	results, total, err := h.service.FetchUserTopTracks(c.Request.Context(), params)
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, err, "Unable to fetch top tracks")
		return
	}

	res := common.BuildPaginatedResponse(results, params.Page, params.Limit, total)

	common.RespondWithData(c, http.StatusOK, res)
}

func (h *ScrobbleStatsHandler) GetUserHistory(c *gin.Context) {
	params, err := h.parsePaginatedStatsParams(c)
	if err != nil {
		common.RespondWithError(c, http.StatusBadRequest, err, "Invalid body")
		return
	}

	results, total, err := h.service.FetchUserHistory(c.Request.Context(), params)
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, err, "Unable to fetch history")
		return
	}

	res := common.BuildPaginatedResponse(results, params.Page, params.Limit, total)

	common.RespondWithData(c, http.StatusOK, res)
}

func (h *ScrobbleStatsHandler) GetArtistTimeline(c *gin.Context) {
	params, err := h.parseStatsParams(c, new(request.ContentTypeArtist))
	if err != nil {
		common.RespondWithError(c, http.StatusBadRequest, err, "Invalid body")
		return
	}

	res, err := h.service.FetchLineChartData(c.Request.Context(), params)
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, err, "Unable to fetch line chart data")
		return
	}

	common.RespondWithData(c, http.StatusOK, res)
}

func (h *ScrobbleStatsHandler) GetAlbumTimeline(c *gin.Context) {
	params, err := h.parseStatsParams(c, new(request.ContentTypeAlbum))
	if err != nil {
		common.RespondWithError(c, http.StatusBadRequest, err, "Invalid body")
		return
	}

	res, err := h.service.FetchLineChartData(c.Request.Context(), params)
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, err, "Unable to fetch line chart data")
		return
	}

	common.RespondWithData(c, http.StatusOK, res)
}

func (h *ScrobbleStatsHandler) GetTrackTimeline(c *gin.Context) {
	params, err := h.parseStatsParams(c, new(request.ContentTypeTrack))
	if err != nil {
		common.RespondWithError(c, http.StatusBadRequest, err, "Invalid body")
		return
	}

	res, err := h.service.FetchLineChartData(c.Request.Context(), params)
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, err, "Unable to fetch line chart data")
		return
	}

	common.RespondWithData(c, http.StatusOK, res)
}
