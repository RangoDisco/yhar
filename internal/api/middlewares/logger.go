package middlewares

import (
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/models"
)

func LoggerMiddleware(l *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := l.With(
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("ip", c.ClientIP()),
		)

		rawUser, exists := c.Get("user")
		if exists && rawUser != nil {
			user := rawUser.(*models.User)
			id := strconv.FormatInt(user.ID, 10)
			logger = logger.With(slog.String("user_id", id))
		}

		c.Set("logger", logger)
		c.Next()
	}
}
