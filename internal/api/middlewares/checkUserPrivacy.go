package middlewares

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/common"
	"github.com/rangodisco/yhar/internal/api/dto"
	"github.com/rangodisco/yhar/internal/api/repositories"
	"github.com/rangodisco/yhar/internal/api/types/filters"
)

func CheckUserPrivacy(repo *repositories.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		rawUser, exists := c.Get("user")
		if !exists {
			common.RespondWithError(c, http.StatusUnauthorized, errors.New("user not found in context"), "Unauthorized")
			return
		}

		currentUser, ok := rawUser.(*dto.UserPassport)
		if !ok && currentUser != nil {
			common.RespondWithError(c, http.StatusInternalServerError, errors.New("unable to convert context user to model"), "Internal server error")
			return
		}

		uID := c.Param("userID")
		if uID == "" {
			common.RespondWithError(c, http.StatusBadRequest, errors.New("userID missing"), "UserID is required")
			return
		}

		// Skip if viewing own data
		if currentUser != nil && (uID == "me" || uID == strconv.FormatInt(currentUser.ID, 10)) {
			c.Next()
			return
		}

		u, err := repo.FindActiveByFilters(ctx, []filters.QueryFilter{
			{Key: "id", Value: uID},
		})

		if err != nil {
			common.RespondWithError(c, http.StatusNotFound, errors.New("user doesn't exist"), "User not found")
			return
		}

		// Only allow public profiles to be seen by other users
		if !u.IsPublic {
			common.RespondWithError(c, http.StatusNotFound, errors.New("user is private"), "User not found")
			return
		}

		c.Next()
	}
}
