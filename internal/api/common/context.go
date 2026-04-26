package common

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/rangodisco/yhar/internal/api/dto"
)

// GetUserFromContext is a helper function that's used in handlers to get the current user from the context
func GetUserFromContext(c *gin.Context) (*dto.UserPassport, error) {
	rawUser, exists := c.Get("user")
	if !exists {
		return nil, errors.New("user not authenticated")
	}
	currentUser, ok := rawUser.(*dto.UserPassport)
	if !ok {
		return nil, errors.New("invalid user")
	}

	return currentUser, nil
}
