package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/services"
	"github.com/rangodisco/yhar/internal/api/types/stats"
	"github.com/rangodisco/yhar/internal/api/utils/convert"
)

type ScrobbleHandler struct {
	scrobbleService *services.ScrobbleService
	statService     *services.ScrobbleStatsService
	subSonicService *services.SubsonicService
}

func NewScrobbleHandler(scrobbleService *services.ScrobbleService, statService *services.ScrobbleStatsService) *ScrobbleHandler {
	return &ScrobbleHandler{scrobbleService: scrobbleService, statService: statService}
}

func parseStatsParams(c *gin.Context) (*stats.Params, error) {
	var userID string
	page := convert.ParseInt(c.Query("page"), 1)
	limit := convert.ParseInt(c.Query("limit"), 10)
	period := stats.Period(c.DefaultQuery("period", "week"))
	paramUserID := c.Param("userID")

	if paramUserID == "me" {
		currentUser, exists := c.Get("user")
		if !exists {
			return nil, errors.New("user not found")
		}
		userID = currentUser.(string)
	} else {
		userID = paramUserID
	}

	return &stats.Params{
		UserID: userID,
		Period: period,
		Pagination: stats.RequestPagination{
			Page:  page,
			Limit: limit,
		},
	}, nil
}

// TODO: handle polling and what is considered a real scrobble
func (h *ScrobbleHandler) ManualNowPlayingPoll(c *gin.Context) {
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
		res, err := h.scrobbleService.HandleNewScrobble(entry)
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
	params, err := parseStatsParams(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	results, total, err := h.statService.FetchUserTopArtists(params)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	res := h.statService.BuildResponseData(results, params.Pagination.Page, params.Pagination.Limit, total)

	c.JSON(http.StatusOK, gin.H{
		"data": res,
	})
}

// GetUserTopAlbums fetches the most scrobbled albums in a given period for a given user
func (h *ScrobbleHandler) GetUserTopAlbums(c *gin.Context) {
	params, err := parseStatsParams(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	results, total, err := h.statService.FetchUserTopAlbums(params)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	res := h.statService.BuildResponseData(results, params.Pagination.Page, params.Pagination.Limit, total)

	c.JSON(http.StatusOK, gin.H{
		"data": res,
	})
}

func (h *ScrobbleHandler) GetUserTopTracks(c *gin.Context) {
	params, err := parseStatsParams(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	results, total, err := h.statService.FetchUserTopTracks(params)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	res := h.statService.BuildResponseData(results, params.Pagination.Page, params.Pagination.Limit, total)

	c.JSON(http.StatusOK, gin.H{
		"data": res,
	})
}
