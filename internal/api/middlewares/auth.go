package middlewares

import "C"
import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/services"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// Proceed as anon if token was not provided
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.Set("user", nil)
			c.Next()
		}

		stringToken := authHeader[7:]
		token, err := services.ParseToken(stringToken)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		// Fetch whole user from token claims
		user, err := services.GetUserFromToken(token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		// Set user to the context
		c.Set("user", user)
		c.Next()
	}
}
