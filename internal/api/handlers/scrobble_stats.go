package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/common"
	"github.com/rangodisco/yhar/internal/api/dto"
	"github.com/rangodisco/yhar/internal/api/repositories"
	"github.com/rangodisco/yhar/internal/api/services"
	"github.com/rangodisco/yhar/internal/api/utils/convert"
)

type ScrobbleStatsHandler struct {
	service *services.ScrobbleStatsService
}

func NewScrobbleStatsHandler(service *services.ScrobbleStatsService) *ScrobbleStatsHandler {
	return &ScrobbleStatsHandler{service: service}
}

// TODO: refacto whole param handling
func (h *ScrobbleStatsHandler) parseStatsParams(c *gin.Context) (*services.PaginatedStatsRequest, error) {
	// Extract user ID
	paramUserID := c.Param("userID")
	var userID string

	if paramUserID == "me" {
		currentUser, err := common.GetUserFromContext(c)
		if err != nil {
			return nil, err
		}
		userID = strconv.Itoa(int(currentUser.ID))
	} else {
		userID = paramUserID
	}

	// Parse and validate pagination
	page := convert.ParseInt(c.Query("page"), 1)
	limit := convert.ParseInt(c.Query("limit"), 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 1
	}
	if limit > 100 {
		limit = 100
	}

	// Parse period
	period := dto.Period(c.DefaultQuery("period", string(dto.PeriodWeek)))

	// Build params
	params := &services.PaginatedStatsRequest{
		UserID: userID,
		Period: period,
		Pagination: struct {
			Page  int
			Limit int
		}{
			Page:  page,
			Limit: limit,
		},
	}

	// Optional filters
	if artistID := c.Query("artist"); artistID != "" {
		params.TargetID = artistID
		params.TargetContent = repositories.Artist
	}

	if trackID := c.Query("track"); trackID != "" {
		params.TargetID = trackID
		params.TargetContent = repositories.Track
	}

	if albumID := c.Query("album"); albumID != "" {
		params.TargetID = albumID
		params.TargetContent = repositories.Album
	}

	return params, nil
}
func (h *ScrobbleStatsHandler) parseRegularStatsParams(c *gin.Context) (*services.StatsRequest, error) {
	// Extract user ID
	paramUserID := c.Param("userID")
	var userID string

	if paramUserID == "me" {
		currentUser, err := common.GetUserFromContext(c)
		if err != nil {
			return nil, err
		}
		userID = strconv.Itoa(int(currentUser.ID))
	} else {
		userID = paramUserID
	}

	// Build params
	params := &services.StatsRequest{
		UserID: userID,
		Period: dto.Period(c.DefaultQuery("period", string(dto.PeriodWeek))),
	}

	// Optional filters
	if artistID := c.Query("artist"); artistID != "" {
		params.TargetID = artistID
		params.TargetContent = repositories.Artist
	}

	if trackID := c.Query("track"); trackID != "" {
		params.TargetID = trackID
		params.TargetContent = repositories.Track
	}

	if albumID := c.Query("album"); albumID != "" {
		params.TargetID = albumID
		params.TargetContent = repositories.Album
	}

	return params, nil
}

// GetUserTopArtists fetches the most scrobbled artists in a given period for a given user
func (h *ScrobbleStatsHandler) GetUserTopArtists(c *gin.Context) {
	ctx := c.Request.Context()
	params, err := h.parseStatsParams(c)
	if err != nil {
		common.RespondWithError(c, http.StatusBadRequest, err, "Invalid body")
		return
	}

	results, total, err := h.service.FetchUserTopArtists(ctx, params)
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, err, "Unable to fetch top artists")
		return
	}

	res := common.BuildPaginatedResponse(results, params.Pagination.Page, params.Pagination.Limit, total)

	common.RespondWithData(c, http.StatusOK, res)
}

// GetUserTopAlbums fetches the most scrobbled albums in a given period for a given user
func (h *ScrobbleStatsHandler) GetUserTopAlbums(c *gin.Context) {
	ctx := c.Request.Context()
	params, err := h.parseStatsParams(c)
	if err != nil {
		common.RespondWithError(c, http.StatusBadRequest, err, "Invalid body")
		return
	}

	results, total, err := h.service.FetchUserTopAlbums(ctx, params)
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, err, "Unable to fetch top albums")
		return
	}

	res := common.BuildPaginatedResponse(results, params.Pagination.Page, params.Pagination.Limit, total)

	common.RespondWithData(c, http.StatusOK, res)
}

func (h *ScrobbleStatsHandler) GetUserTopTracks(c *gin.Context) {
	ctx := c.Request.Context()
	params, err := h.parseStatsParams(c)
	if err != nil {
		common.RespondWithError(c, http.StatusBadRequest, err, "Invalid body")
		return
	}

	results, total, err := h.service.FetchUserTopTracks(ctx, params)
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, err, "Unable to fetch top tracks")
		return
	}

	res := common.BuildPaginatedResponse(results, params.Pagination.Page, params.Pagination.Limit, total)

	common.RespondWithData(c, http.StatusOK, res)
}

func (h *ScrobbleStatsHandler) GetUserHistory(c *gin.Context) {
	ctx := c.Request.Context()

	params, err := h.parseStatsParams(c)
	if err != nil {
		common.RespondWithError(c, http.StatusBadRequest, err, "Invalid body")
		return
	}

	results, total, err := h.service.FetchUserHistory(ctx, params)
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, err, "Unable to fetch history")
		return
	}

	res := common.BuildPaginatedResponse(results, params.Pagination.Page, params.Pagination.Limit, total)

	common.RespondWithData(c, http.StatusOK, res)
}

func (h *ScrobbleStatsHandler) GetArtistTimeline(c *gin.Context) {
	ctx := c.Request.Context()
	params := services.StatsRequest{
		UserID: c.Param("userID"),
		Period: dto.Period(c.DefaultQuery("period", string(dto.PeriodWeek))),
	}

	artistID := c.Param("artistID")
	if artistID == "" {
		common.RespondWithError(c, http.StatusInternalServerError, errors.New("no artistID provided"), "Unable to fetch line chart data")
		return
	}

	params.TargetContent = repositories.Artist
	params.TargetID = artistID
	params.End = new(time.Now())
	params.Start = new(time.Now().Add(time.Duration(-48) * time.Hour))

	res, err := h.service.FetchLineChartData(ctx, &params)
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, err, "Unable to fetch line chart data")
		return
	}

	common.RespondWithData(c, http.StatusOK, res)
}

func (h *ScrobbleStatsHandler) GetAlbumTimeline(c *gin.Context) {
	ctx := c.Request.Context()
	params := services.StatsRequest{
		UserID: c.Param("userID"),
		Period: dto.Period(c.DefaultQuery("period", string(dto.PeriodWeek))),
	}

	albumID := c.Param("albumID")
	if albumID == "" {
		common.RespondWithError(c, http.StatusInternalServerError, errors.New("to albumID provided"), "Unable to fetch line chart data")
		return
	}

	params.TargetContent = repositories.Album
	params.TargetID = albumID
	params.End = new(time.Now())
	params.Start = new(time.Now().Add(time.Duration(-48) * time.Hour))

	res, err := h.service.FetchLineChartData(ctx, &params)
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, err, "Unable to fetch line chart data")
		return
	}

	common.RespondWithData(c, http.StatusOK, res)
}

func (h *ScrobbleStatsHandler) GetTrackTimeline(c *gin.Context) {
	ctx := c.Request.Context()
	params := services.StatsRequest{
		UserID: c.Param("userID"),
		Period: dto.Period(c.DefaultQuery("period", string(dto.PeriodWeek))),
	}

	trackID := c.Param("trackID")
	if trackID == "" {
		common.RespondWithError(c, http.StatusInternalServerError, errors.New("no trackID provided"), "Unable to fetch line chart data")
		return
	}

	params.TargetContent = repositories.Track
	params.TargetID = trackID
	params.End = new(time.Now())
	params.Start = new(time.Now().Add(time.Duration(-48) * time.Hour))

	res, err := h.service.FetchLineChartData(ctx, &params)
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, err, "Unable to fetch line chart data")
		return
	}

	common.RespondWithData(c, http.StatusOK, res)
}
