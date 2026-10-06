package common

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// ResolveUserID extracts the user ID from the "userID" URL parameter.
func ResolveUserID(c *gin.Context) (string, error) {
	paramUserID := c.Param("userID")
	if paramUserID == "me" {
		currentUser, err := GetUserFromContext(c)
		if err != nil {
			return "", err
		}
		return strconv.Itoa(int(currentUser.ID)), nil
	}
	return paramUserID, nil
}
