package middlewares

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/common"
	"github.com/rangodisco/yhar/internal/api/dto"
	"github.com/rangodisco/yhar/internal/api/services"
)

func CheckUserPrivacy(service *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		rawUser, exists := c.Get("user")
		if !exists {
			common.RespondWithError(c, http.StatusUnauthorized, errors.New("user not found in context"), "Unauthorized")
			return
		}

		currentUser, ok := rawUser.(*dto.UserPassport)
		if rawUser != nil && !ok {
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

		iID, err := strconv.ParseInt(uID, 10, 64)
		if err != nil {
			common.RespondWithError(c, http.StatusBadRequest, errors.New("invalid userID"), "Invalid userID")
			return
		}

		u, err := service.GetUserByID(ctx, iID)
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
