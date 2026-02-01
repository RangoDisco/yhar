package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/services"
	"github.com/rangodisco/yhar/internal/api/types/stats"
	"github.com/rangodisco/yhar/internal/api/utils/convert"
)

func parseStatsParams(c *gin.Context) (*stats.Params, error) {
	page := convert.ParseInt(c.Query("page"), 1)
	limit := convert.ParseInt(c.Query("limit"), 10)
	period := stats.Period(c.DefaultQuery("period", "week"))
	userID := c.Param("userID")

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
func ManualNowPlayingPoll(c *gin.Context) {
	var scrobbles []*models.Scrobble
	user, exists := c.Get("user")
	if !exists {
		c.JSON(500, gin.H{})
		return
	}
	fmt.Println(user)

	subRes, err := services.GetNowPlaying()
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
		res, err := services.HandleNewScrobble(entry)
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
func GetUserTopArtists(c *gin.Context) {
	params, err := parseStatsParams(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	results, total, err := services.FetchUserTopArtists(params)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	res := services.BuildResponseData(results, params.Pagination.Page, params.Pagination.Limit, total)

	c.JSON(http.StatusOK, gin.H{
		"data": res,
	})
}

// GetUserTopAlbums fetches the most scrobbled albums in a given period for a given user
func GetUserTopAlbums(c *gin.Context) {
	params, err := parseStatsParams(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	results, total, err := services.FetchUserTopAlbums(params)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	res := services.BuildResponseData(results, params.Pagination.Page, params.Pagination.Limit, total)

	c.JSON(http.StatusOK, gin.H{
		"data": res,
	})
}

func GetUserTopTracks(c *gin.Context) {
	params, err := parseStatsParams(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	results, total, err := services.FetchUserTopTracks(params)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	res := services.BuildResponseData(results, params.Pagination.Page, params.Pagination.Limit, total)

	c.JSON(http.StatusOK, gin.H{
		"data": res,
	})
}
