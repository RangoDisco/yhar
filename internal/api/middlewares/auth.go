package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/common"
	"github.com/rangodisco/yhar/internal/api/services"
)

func Authenticate(auth *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		authHeader := c.GetHeader("Authorization")

		// Proceed as anon if token was not provided
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.Set("user", nil)
			c.Next()
		} else {
			stringToken := authHeader[7:]
			token, err := services.ParseToken(stringToken)
			if err != nil {
				common.RespondWithError(c, http.StatusUnauthorized, err, "Unauthorized")
				return
			}

			// Fetch whole user from token claims
			user, err := auth.GetUserFromToken(ctx, token)
			if err != nil {
				common.RespondWithError(c, http.StatusUnauthorized, err, "Unauthorized")
				return
			}

			// TODO: don't set the whole user
			// Set user to the context
			c.Set("user", user)
			c.Next()
		}
	}
}
