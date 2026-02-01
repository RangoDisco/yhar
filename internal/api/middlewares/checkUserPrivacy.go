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
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		currentUser, ok := rawUser.(*models.User)
		if !ok {
			c.JSON(500, gin.H{"error": "server error"})
			c.Abort()
			return
		}

		uID := c.Param("userID")
		if uID == "" {
			c.JSON(400, gin.H{"error": "userID param is missing"})
			c.Abort()
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
			c.JSON(404, gin.H{"error": "user not found"})
			c.Abort()
			return
		}

		// Only allow public profiles to be seen by other users
		if !u.IsPublic {
			c.JSON(404, gin.H{"error": "not found"})
			c.Abort()
			return
		}

		c.Next()
	}
}
