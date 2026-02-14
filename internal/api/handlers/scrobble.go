package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/dto"
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/services"
	"github.com/rangodisco/yhar/internal/api/utils/convert"
)

type ScrobbleHandler struct {
	scrobbleService *services.ScrobbleService
	statService     *services.ScrobbleStatsService
	subSonicService *services.SubsonicService
}

type QueryParams struct {
	Period dto.Period
	Page   int
	Limit  int
}

func NewScrobbleHandler(scrobbleService *services.ScrobbleService, statService *services.ScrobbleStatsService) *ScrobbleHandler {
	return &ScrobbleHandler{scrobbleService: scrobbleService, statService: statService}
}

func (h *ScrobbleHandler) parseStatsParams(c *gin.Context) (*services.StatsRequestParams, error) {
	// Extract user ID
	paramUserID := c.Param("userID")
	var userID string

	if paramUserID == "me" {
		rawUser, exists := c.Get("user")
		if !exists {
			return nil, errors.New("user not authenticated")
		}
		currentUser, ok := rawUser.(*models.User)
		if !ok {
			return nil, errors.New("invalid user")
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
	params := &services.StatsRequestParams{
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
		params.ArtistID = &artistID
	}

	if trackID := c.Query("track"); trackID != "" {
		params.TrackID = &trackID
	}

	return params, nil
}

// TODO: handle polling and what is considered a real scrobble
func (h *ScrobbleHandler) ManualNowPlayingPoll(c *gin.Context) {
	ctx := c.Request.Context()
	var scrobbles []*models.Scrobble
	user, exists := c.Get("user")
	if !exists {
		c.JSON(500, gin.H{})
		return
	}
	fmt.Println(user)

	subRes, err := h.subSonicService.GetNowPlaying()
	if err != nil {
		c.JSON(500, gin.H{})
		return
	}

	if len(subRes.NowPlaying.Entry) == 0 {
		c.JSON(200, gin.H{})
		return
	}

	for _, entry := range subRes.NowPlaying.Entry {
		// Check if track already exists in service's db
		res, err := h.scrobbleService.HandleNewScrobble(ctx, entry)
		if err != nil {
			continue
		}
		scrobbles = append(scrobbles, res)
	}

	c.JSON(200, gin.H{
		"data": scrobbles,
	})
}

// GetUserTopArtists fetches the most scrobbled artists in a given period for a given user
func (h *ScrobbleHandler) GetUserTopArtists(c *gin.Context) {
	ctx := c.Request.Context()
	params, err := h.parseStatsParams(c)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, err, "Invalid body")
		return
	}

	results, total, err := h.statService.FetchUserTopArtists(ctx, params)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, err, "Unable to fetch top artists")
		return
	}

	res := BuildPaginatedResponse(results, params.Pagination.Page, params.Pagination.Limit, total)

	RespondWithData(c, http.StatusOK, res)
}

// GetUserTopAlbums fetches the most scrobbled albums in a given period for a given user
func (h *ScrobbleHandler) GetUserTopAlbums(c *gin.Context) {
	ctx := c.Request.Context()
	params, err := h.parseStatsParams(c)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, err, "Invalid body")
		return
	}

	results, total, err := h.statService.FetchUserTopAlbums(ctx, params)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, err, "Unable to fetch top albums")
		return
	}

	res := BuildPaginatedResponse(results, params.Pagination.Page, params.Pagination.Limit, total)

	RespondWithData(c, http.StatusOK, res)
}

func (h *ScrobbleHandler) GetUserTopTracks(c *gin.Context) {
	ctx := c.Request.Context()
	params, err := h.parseStatsParams(c)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, err, "Invalid body")
		return
	}

	results, total, err := h.statService.FetchUserTopTracks(ctx, params)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, err, "Unable to fetch top tracks")
		return
	}

	res := BuildPaginatedResponse(results, params.Pagination.Page, params.Pagination.Limit, total)

	RespondWithData(c, http.StatusOK, res)
}

func (h *ScrobbleHandler) GetUserHistory(c *gin.Context) {
	ctx := c.Request.Context()

	params, err := h.parseStatsParams(c)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, err, "Invalid body")
		return
	}

	results, total, err := h.statService.FetchUserHistory(ctx, params)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, err, "Unable to fetch history")
		return
	}

	res := BuildPaginatedResponse(results, params.Pagination.Page, params.Pagination.Limit, total)

	RespondWithData(c, http.StatusOK, res)
}
