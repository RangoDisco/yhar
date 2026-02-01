package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/models"
)

func RequirePermissions(perms []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawUser, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
		}

		u := rawUser.(*models.User)
		uPerms := make(map[string]bool)
		for _, p := range u.Role.Permissions {
			uPerms[p.Name] = true
		}

		authorized := false
		for _, required := range perms {
			if uPerms[required] {
				authorized = true
				break
			}
		}

		if !authorized {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}
