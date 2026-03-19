package middlewares

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/common"
	"github.com/rangodisco/yhar/internal/api/dto"
)

func RequirePermissions(perms []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawUser, exists := c.Get("user")
		if !exists {
			common.RespondWithError(c, http.StatusForbidden, errors.New("no user found in context"), "Unauthorized")
			return
		}

		u := rawUser.(*dto.UserPassport)
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
			common.RespondWithError(c, http.StatusForbidden, errors.New("user doesn't have the required permissions"), "Forbidden")
			return
		}
		c.Next()
	}
}
