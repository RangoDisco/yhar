package handlers

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

type APIResponse[T any] struct {
	Data T `json:"data,omitempty"`
}

type ErrorResponse struct {
	Error ErrorMessage `json:"error"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

type PaginatedResponse[T any] struct {
	Results    T          `json:"results"`
	Pagination Pagination `json:"pagination"`
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

	c.AbortWithStatusJSON(statusCode, ErrorResponse{Error: ErrorMessage{
		Message: message,
	}})
}

func RespondWithData[T any](c *gin.Context, statusCode int, data T) {
	c.JSON(statusCode, APIResponse[T]{Data: data})
}

func BuildPaginatedResponse[T any](results T, page, limit int, totalCount int64) PaginatedResponse[T] {
	totalPages := int(totalCount) / limit
	if int(totalCount)%limit > 0 {
		totalPages++
	}

	return PaginatedResponse[T]{
		Results: results,
		Pagination: Pagination{
			TotalCount:  totalCount,
			HasNextPage: page < totalPages,
		},
	}
}
