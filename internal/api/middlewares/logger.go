package middlewares

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/dto"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		args := []any{
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("ip", c.ClientIP()),
			slog.String("duration", time.Since(startTime).String()),
			slog.Int("status", c.Writer.Status()),
		}

		// User might not be set, in case of a public route
		rawUser, exists := c.Get("user")
		if exists && rawUser != nil {
			user, ok := rawUser.(*dto.UserPassport)
			if ok {
				id := strconv.FormatInt(user.ID, 10)
				args = append(args, slog.String("user_id", id))
			}
		}

		if len(c.Errors) > 0 {
			args = append(args, slog.String("error", c.Errors.Last().Error()))
			slog.Error("Request completed with errors", args...)
		} else {
			slog.Info("Request completed", args...)
		}
	}
}
