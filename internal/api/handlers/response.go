package handlers

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

type PaginatedResponse struct {
	Results    interface{} `json:"results"`
	Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
	TotalCount  int64 `json:"total_count"`
	HasNextPage bool  `json:"has_next_page"`
}

func RespondWithError(c *gin.Context, statusCode int, err error, message string) {
	rawLogger, exists := c.Get("logger")
	logger, ok := rawLogger.(*slog.Logger)
	if exists && ok {
		logger.Error(message,
			slog.Int("status_code", statusCode),
			slog.String("error", err.Error()),
		)
	}

	c.AbortWithStatusJSON(statusCode, gin.H{"error": message})
}

func RespondWithData(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{"data": data})
}

func BuildPaginatedResponse(results interface{}, page, limit int, totalCount int64) interface{} {
	totalPages := int(totalCount) / limit
	if int(totalCount)%limit > 0 {
		totalPages++
	}

	return &PaginatedResponse{
		Results: results,
		Pagination: Pagination{
			TotalCount:  totalCount,
			HasNextPage: page < totalPages,
		},
	}
}
