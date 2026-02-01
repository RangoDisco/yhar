package middlewares

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/repositories"
	"github.com/rangodisco/yhar/internal/api/types/filters"
)

func CheckUserPrivacy() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawUser, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}

		currentUser, ok := rawUser.(*models.User)
		if !ok {
			c.AbortWithStatusJSON(500, gin.H{"error": "server error"})
			return
		}

		uID := c.Param("userID")
		if uID == "" {
			c.AbortWithStatusJSON(400, gin.H{"error": "userID param is missing"})
			return
		}

		// Skip if viewing own data
		if uID == strconv.FormatInt(currentUser.ID, 10) {
			c.Next()
			return
		}

		rFilters := []filters.QueryFilter{
			{Key: "id", Value: uID},
		}

		u, err := repositories.FindActiveUserByFilters(rFilters)
		if err != nil {
			c.AbortWithStatusJSON(404, gin.H{"error": "user not found"})
			return
		}

		// Only allow public profiles to be seen by other users
		if !u.IsPublic {
			c.AbortWithStatusJSON(404, gin.H{"error": "not found"})
			return
		}

		c.Next()
	}
}
