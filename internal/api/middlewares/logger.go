package middlewares

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
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
			slog.String("user_agent", c.Request.UserAgent()),
			slog.String("duration", time.Since(startTime).String()),
			slog.Int("status", c.Writer.Status()),
		}

		hub := sentrygin.GetHubFromContext(c)

		// User might not be set, in case of a public route
		rawUser, exists := c.Get("user")
		if exists && rawUser != nil {
			user, ok := rawUser.(*dto.UserPassport)
			if ok {
				id := strconv.FormatInt(user.ID, 10)
				args = append(args, slog.String("user_id", id))

				if hub != nil {
					hub.Scope().SetUser(sentry.User{ID: id})
				}
			}
		}

		if len(c.Errors) > 0 {
			last := c.Errors.Last()
			args = append(args, slog.Any("error", last))
			slog.ErrorContext(c.Request.Context(), "Request completed with errors", args...)

			if hub != nil {
				hub.CaptureException(last.Err)
			}
		} else {
			slog.InfoContext(c.Request.Context(), "Request completed", args...)
		}
	}
}
